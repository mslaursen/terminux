package terminux

type cell struct {
	character rune
	color     string
}

func newCell(c rune, col string) *cell {
	return &cell{c, col}
}

func newCellEmpty() *cell {
	return &cell{
		character: ' ',
	}
}

func attachColor(c rune, col string) string {
	return string(c) + string(col)
}

func newBuffer2D[T any](w, h int) [][]T {
	buf := make([][]T, h)
	for i := range h {
		buf[i] = make([]T, w)
	}
	return buf
}
