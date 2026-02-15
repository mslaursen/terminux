package terminux

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type ScreenConfig struct {
	Width, Height int
}

type Screen struct {
	width, height  int
	fd             int
	inputReader    *bufio.Reader
	outWriter      *bufio.Writer
	oldState       *term.State
	currCellBuffer *cellBuffer
	prevCellBuffer *cellBuffer
	lineBuilder    *strings.Builder
	eventListener  func(*Event)
	events         chan *Event
}

func NewScreenDefault() *Screen {
	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		log.Fatal(err)
	}
	return NewScreen(&ScreenConfig{Width: width, Height: height})
}

func NewScreen(cfg *ScreenConfig) *Screen {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		log.Fatal("stdin is not a terminal")
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}

	initEventRegexMap()

	return &Screen{
		inputReader:    bufio.NewReader(os.Stdin),
		outWriter:      bufio.NewWriter(os.Stdout),
		oldState:       oldState,
		fd:             fd,
		currCellBuffer: newCellBuffer(cfg.Width, cfg.Height),
		prevCellBuffer: newCellBuffer(cfg.Width, cfg.Height),
		events:         make(chan *Event, 64),
		width:          cfg.Width,
		height:         cfg.Height,
		lineBuilder:    &strings.Builder{},
	}
}

func (s *Screen) Clear() {
	s.currCellBuffer.clear()
}

// diff check + flush buffer to stdout
func (s *Screen) Display() {
	s.lineBuilder.Reset()
	s.lineBuilder.Grow(s.width * s.height * 16)

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			curr := s.currCellBuffer.getCell(x, y)
			prev := s.prevCellBuffer.getCell(x, y)
			if curr.equal(prev) {
				continue
			}
			s.lineBuilder.WriteString(s.getCursorString(x, y))
			s.lineBuilder.WriteString(
				attachColor(curr.char, curr.fg, curr.bg),
			)
			s.prevCellBuffer.setCell(x, y, curr)
		}
	}

	s.outWriter.WriteString(s.lineBuilder.String())
	s.outWriter.Flush()
}

func (s *Screen) moveCursor(x, y int) {
	fmt.Fprintf(s.outWriter, string(ansiCursorMove), y+1, x+1)
}

func (s *Screen) getCursorString(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

func (s *Screen) Debug(val any, x, y int) {
	s.moveCursor(x, y)
	fmt.Fprintf(s.outWriter, "%v", val)
}

func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

func (s *Screen) Restore() {
	s.outWriter.WriteString(string(ansiShowCursor))
	s.outWriter.WriteString(string(ansiClearScreen))
	s.outWriter.WriteString(string(ansiCursorHome))
	s.outWriter.WriteString(string(ansiDisableMouse))
	s.outWriter.WriteString(string(ansiDisableMouseSGR))
	s.outWriter.Flush()
	term.Restore(s.fd, s.oldState)
}

func (s *Screen) EnableMouse() {
	fmt.Print(ansiEnableMouse)
	fmt.Print(ansiEnableMouseSGR)
}

func (s *Screen) HideCursor() {
	fmt.Print(ansiHideCursor)
}

func (s *Screen) Ticker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
