package terminux

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/term"
)

type EventType int

const (
	AnsiEsc             = "\x1b"
	AnsiCursorMove      = "\033[%d;%dH"
	AnsiHideCursor      = "\033[?25l"
	AnsiShowCursor      = "\033[?25h"
	AnsiClearScreen     = "\033[2J"
	AnsiCursorHome      = "\033[H"
	AnsiEnableMouse     = "\x1b[?1000h"
	AnsiDisableMouse    = "\x1b[?1000l"
	AnsiEnableMouseSGR  = "\x1b[?1006h"
	AnsiDisableMouseSGR = "\x1b[?1006l"
)

const (
	MousePressed EventType = iota
	MouseReleased
	KeyPressed
	Error
)

type Event struct {
	Type EventType
	X, Y int
	Key  string
}

func NewEventError() *Event {
	return NewEvent(Error, -1, -1, "")
}

func NewEvent(t EventType, x, y int, k string) *Event {
	return &Event{
		Type: t,
		X:    x,
		Y:    y,
		Key:  k,
	}
}

var EventRegexMap = make(map[EventType]string)

func InitEventRegexMap() {
	EventRegexMap[MousePressed] = `\x1b\[<0;([0-9]+);([0-9]+)m`
	EventRegexMap[MouseReleased] = `\x1b\[<0;([0-9]+);([0-9]+)M`
}

type ScreenConfig struct {
	Width, Height int
}

type Screen struct {
	width, height int
	fd            int
	inputReader   *bufio.Reader
	outWriter     *bufio.Writer
	oldState      *term.State
	backBuffer    [][]rune
	frontBuffer   [][]rune
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

	front := make([][]rune, cfg.Height)
	back := make([][]rune, cfg.Height)

	for i := range cfg.Height {
		front[i] = make([]rune, cfg.Width)
		back[i] = make([]rune, cfg.Width)
	}

	InitEventRegexMap()

	return &Screen{
		inputReader: bufio.NewReader(os.Stdin),
		outWriter:   bufio.NewWriter(os.Stdout),
		oldState:    oldState,
		fd:          fd,
		backBuffer:  back,
		frontBuffer: front,
		events:      make(chan *Event, 64),
		width:       cfg.Width,
		height:      cfg.Height,
	}
}

func (s *Screen) Clear() {
	for y := range s.height {
		for x := range s.width {
			s.backBuffer[y][x] = ' '
		}
	}
}

func (s *Screen) Draw(x, y int, c rune) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return
	}
	s.backBuffer[y][x] = c
}

func (s *Screen) DrawRect(x, y, w, h int, fill bool, c rune) {
	for dx := range w {
		for dy := range h {
			if !fill {
				if dx != 0 && dx != w-1 && dy != 0 && dy != h-1 {
					continue
				}
			}
			s.Draw(dx+x, dy+y, c)
		}
	}
}
func (s *Screen) DrawLine(x0, y0, x1, y1 int, c rune) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx := -1
	if x0 < x1 {
		sx = 1
	}

	sy := -1
	if y0 < y1 {
		sy = 1
	}

	err := dx - dy

	for {
		s.Draw(x0, y0, c)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// diff check + flush buffer to stdout
func (s *Screen) Display() {
	for y := range s.height {
		for x := range s.width {
			curr := s.backBuffer[y][x]
			prev := s.frontBuffer[y][x]
			if curr != prev {
				fmt.Fprintf(s.outWriter, AnsiCursorMove, y+1, x+1)
				s.outWriter.WriteRune(curr)
				s.frontBuffer[y][x] = curr
			}
		}
	}
	s.outWriter.Flush()
}

func (s *Screen) Debug(val any, x, y int) {
	fmt.Fprintf(s.outWriter, AnsiCursorMove, y+1, x+1)
	fmt.Fprintf(s.outWriter, "%v", val)
}

func (s *Screen) ListenForEvents() {
	buf := make([]byte, 64)

	for {
		n, _ := s.inputReader.Read(buf)
		if buf[0] == 0x1b {
			ev, err := s.ParseANSIEvent(buf[:n])
			if err != nil {
				ev = NewEventError()
				s.Debug(err, 0, 1)
			}
			s.events <- ev
			continue
		}
		s.events <- NewEvent(KeyPressed, -1, -1, string(buf[0]))
	}
}

func (s *Screen) Events() <-chan *Event {
	return s.events
}

func (s *Screen) ParseANSIEvent(b []byte) (*Event, error) {
	v := string(b)
	for t, r := range EventRegexMap {
		re := regexp.MustCompile(r)
		matches := re.FindStringSubmatch(v)
		if len(matches) == 3 {
			x, _ := strconv.Atoi(matches[1])
			y, _ := strconv.Atoi(matches[2])
			return &Event{Type: t, X: x, Y: y}, nil
		}
	}
	return nil, errors.New("failed to parse ANSI string")
}

func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

func (s *Screen) Restore() {
	s.outWriter.WriteString(AnsiShowCursor)
	s.outWriter.WriteString(AnsiClearScreen)
	s.outWriter.WriteString(AnsiCursorHome)
	s.outWriter.WriteString(AnsiDisableMouse)
	s.outWriter.WriteString(AnsiDisableMouseSGR)
	s.outWriter.Flush()
	term.Restore(s.fd, s.oldState)
}

func (s *Screen) EnableMouse() {
	fmt.Print(AnsiEnableMouse)
	fmt.Print(AnsiEnableMouseSGR)
}

func (s *Screen) HideCursor() {
	fmt.Print(AnsiHideCursor)
}

func (s *Screen) Ticker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
