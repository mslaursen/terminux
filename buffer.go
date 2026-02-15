package terminux

import "fmt"

type cell struct {
	char rune
	fg   uint32
	bg   uint32
}

type cellBuffer struct {
	cells         [][]*cell
	width, height int
}

func newCell(char rune, fg, bg uint32) *cell {
	return &cell{char, fg, bg}
}

func newCellEmpty() *cell {
	return &cell{
		char: ' ',
	}
}

func newCellBuffer(w, h int) *cellBuffer {
	buf := make([][]*cell, h)
	for i := range h {
		buf[i] = make([]*cell, w)
	}
	return &cellBuffer{
		buf, w, h,
	}

}

func (cb *cellBuffer) clear() {
	for y := range cb.height {
		for x := range cb.width {
			cb.cells[y][x] = newCellEmpty()
		}
	}
}

func (cb *cellBuffer) getCell(x, y int) *cell {
	return cb.cells[y][x]
}

func (cb *cellBuffer) setCell(x, y int, cell *cell) {
	cb.cells[y][x] = cell
}

func (cb *cell) equal(other *cell) bool {
	if other == nil {
		return false
	}
	return cb.char == other.char &&
		cb.fg == other.fg &&
		cb.bg == other.bg
}

func attachColor(char rune, styles ...uint32) string {
	styleString := ""
	for i, s := range styles {
		if i > 0 {
			styleString += ";"
		}
		styleString += fmt.Sprintf("%d", s)
	}

	return fmt.Sprintf("\033[%sm%c\033[0m", styleString, char)
}
