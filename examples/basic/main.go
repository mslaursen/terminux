package main

import (
	"math"
	"time"

	"github.com/mslaursen/terminux"
)

func main() {
	screen := terminux.NewScreenDefault()
	screen.HideCursor()
	screen.EnableMouse()
	defer screen.Restore()

	go screen.ListenForEvents()

	ticker := screen.Ticker(time.Second / 30)
	defer ticker.Stop()

	t := 0.0

	x0, y0 := 0, 0
	x1, y1 := 0, 0
	for {
		select {
		case ev := <-screen.Events():
			if ev.Type == terminux.KeyPressed && ev.Key == "q" {
				return
			}
			if ev.Type == terminux.MousePressed {
				x0 = ev.X
				y0 = ev.Y
			}
			if ev.Type == terminux.MouseReleased {
				x1 = ev.X
				y1 = ev.Y
			}

		case dt := <-ticker.C:
			screen.Debug(dt, 0, 0)
			t += 0.1
			screen.Clear()
			screen.DrawRect(int(math.Cos(t)*10)+50, 10, 5, 5, true, '#')
			screen.DrawLine(x0, y0, x1, y1, '@')
			screen.Display()
		}
	}
}
