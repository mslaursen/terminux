package terminux

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"
)

type ScreenConfig struct {
	Width, Height int
}

type Screen struct {
	width, height int
	fd            int
	inputReader   *bufio.Reader
	outWriter     *bufio.Writer
	oldState      *term.State
	backBuffer    [][]*cell
	frontBuffer   [][]*cell
	eventListener func(*Event)
	events        chan *Event
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
		inputReader: bufio.NewReader(os.Stdin),
		outWriter:   bufio.NewWriter(os.Stdout),
		oldState:    oldState,
		fd:          fd,
		frontBuffer: newBuffer2D[*cell](cfg.Width, cfg.Height),
		backBuffer:  newBuffer2D[*cell](cfg.Width, cfg.Height),
		events:      make(chan *Event, 64),
		width:       cfg.Width,
		height:      cfg.Height,
	}
}

// diff check + flush buffer to stdout
func (s *Screen) Display() {
	for y := range s.height {
		for x := range s.width {
			curr := s.backBuffer[y][x]
			prev := s.frontBuffer[y][x]
			if curr != prev {
				fmt.Fprintf(s.outWriter, ansiCursorMove, y+1, x+1)
				s.outWriter.WriteString(attachColor(curr.character, curr.color))
				s.frontBuffer[y][x] = curr
			}
		}
	}
	s.outWriter.Flush()
}

func (s *Screen) Debug(val any, x, y int) {
	fmt.Fprintf(s.outWriter, ansiCursorMove, y+1, x+1)
	fmt.Fprintf(s.outWriter, "%v", val)
}

func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

func (s *Screen) Restore() {
	s.outWriter.WriteString(ansiShowCursor)
	s.outWriter.WriteString(ansiClearScreen)
	s.outWriter.WriteString(ansiCursorHome)
	s.outWriter.WriteString(ansiDisableMouse)
	s.outWriter.WriteString(ansiDisableMouseSGR)
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
