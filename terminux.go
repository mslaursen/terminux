package terminux

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"

	"golang.org/x/term"
)

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
	eventListener func(string)
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
	return &Screen{
		inputReader:   bufio.NewReader(os.Stdin),
		outWriter:     bufio.NewWriter(os.Stdout),
		oldState:      oldState,
		fd:            fd,
		eventListener: nil,
		backBuffer:    back,
		frontBuffer:   front,
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
	s.frontBuffer, s.backBuffer = s.backBuffer, s.frontBuffer
	s.outWriter.Flush()
}

func (s *Screen) SetEventListener(eventListener func(string)) {
	s.eventListener = eventListener
}

func (s *Screen) ListenForEvents() {
	for {
		s.eventListener(s.getEvent())
	}
}

func (s *Screen) getEvent() string {
	r, _, _ := s.inputReader.ReadRune()
	return string(r)
}

func (s *Screen) GetSize() (int, int) {
	return s.Width, s.Height
}

func (s *Screen) Restore() {
	s.outWriter.WriteString("\033[?25h")
	s.outWriter.WriteString("\033[2J\033[H")
	s.outWriter.Flush()
	term.Restore(s.fd, s.oldState)
}

func (s *Screen) EnableMouse() {
	fmt.Print("\x1b[?1000h")
}

func (s *Screen) HideCursor() {
	fmt.Print("\033[?25l")
}
