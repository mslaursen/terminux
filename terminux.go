package terminux

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/term"
)

type EventType int

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
	Width, Height int
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
		inputReader:   bufio.NewReader(os.Stdin),
		outWriter:     bufio.NewWriter(os.Stdout),
		oldState:      oldState,
		fd:            fd,
		eventListener: nil,
		backBuffer:    back,
		frontBuffer:   front,
		events:        make(chan *Event, 64),
		Width:         cfg.Width,
		Height:        cfg.Height,
	}
}

func (s *Screen) Clear() {
	for y := range s.Height {
		for x := range s.Width {
			s.backBuffer[y][x] = ' '
		}
	}
}

func (s *Screen) Draw(x, y int, c rune) {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
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

func (s *Screen) DrawLine(x1, y1, x2, y2 int, c rune) {
	vx, vy := float64(x2-x1), float64(y2-y1)
	vl := math.Hypot(vx, vy)
	if vl == 0 {
		s.Draw(int(vx), int(vy), c)
	}
	vnx, vny := vx/vl, vy/vl
	for step := range int(vl) {
		fs := float64(step)
		x := x1 + int(math.Round(vnx*fs))
		y := y1 + int(math.Round(vny*fs))
		s.Draw(x, y, c)
	}
}

// diff check + flush buffer to stdout
func (s *Screen) Display() {
	for y := range s.Height {
		for x := range s.Width {
			curr := s.backBuffer[y][x]
			prev := s.frontBuffer[y][x]
			if curr != prev {
				fmt.Fprintf(s.outWriter, "\033[%d;%dH%c", y+1, x+1, curr)
				s.frontBuffer[y][x] = curr
			}
		}
	}
	s.outWriter.Flush()
}

func (s *Screen) Debug(val any, x, y int) {
	fmt.Fprintf(s.outWriter, "\033[%d;%dH%s", y+1, x+1, fmt.Sprintf("%v", val))
}

func (s *Screen) SetEventListener(eventListener func(*Event)) {
	s.eventListener = eventListener
}

func (s *Screen) ListenForEvents() {
	buf := make([]byte, 64)

	for {
		n, _ := s.inputReader.Read(buf)
		if buf[0] == 0x1b {
			ev, err := s.ParseANSIEvent(buf[:n])
			if err != nil {
				ev = NewEventError()
				s.Debug(err, 10, 10)
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
			event := &Event{
				Type: t,
				X:    x,
				Y:    y,
			}
			return event, nil
		}
	}
	return nil, errors.New("failed to parse ANSI string")
}

func (s *Screen) Size() (int, int) {
	return s.Width, s.Height
}

func (s *Screen) Restore() {
	s.outWriter.WriteString("\033[?25h")
	s.outWriter.WriteString("\033[2J\033[H")
	fmt.Print("\x1b[?1000l")
	fmt.Print("\x1b[?1006l")
	s.outWriter.Flush()
	term.Restore(s.fd, s.oldState)
}

func (s *Screen) EnableMouse() {
	fmt.Print("\x1b[?1000h")
	fmt.Print("\x1b[?1006h")

}

func (s *Screen) HideCursor() {
	fmt.Print("\033[?25l")
}

func (s *Screen) Ticker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
