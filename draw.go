package terminux

import (
	"errors"
	"regexp"
	"strconv"
)

func (s *Screen) Draw(x, y int, c rune, fg, bg uint32) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return
	}
	s.currCellBuffer.setCell(x, y, newCell(c, fg, bg))
}

func (s *Screen) DrawRect(x, y, w, h int, fill bool, c rune, fg, bg uint32) {
	for dx := range w {
		for dy := range h {
			if !fill {
				if dx != 0 && dx != w-1 && dy != 0 && dy != h-1 {
					continue
				}
			}
			s.Draw(dx+x, dy+y, c, fg, bg)
		}
	}
}

func (s *Screen) DrawLine(x0, y0, x1, y1 int, c rune, fg, bg uint32) {
	dx := intAbs(x1 - x0)
	dy := intAbs(y1 - y0)

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
		s.Draw(x0, y0, c, fg, bg)

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

func intAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func parseANSIEvent(b []byte) (*Event, error) {
	v := string(b)
	for t, r := range eventRegexMap {
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
