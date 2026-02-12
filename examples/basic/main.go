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

	x1, y1 := 0, 0
	x2, y2 := 0, 0
	for {
		select {
		case ev := <-screen.Events():
			if ev.Type == terminux.KeyPressed && ev.Key == "q" {
				return
			}
			if ev.Type == terminux.MousePressed {
				x1 = ev.X
				y1 = ev.Y
			}
			if ev.Type == terminux.MouseReleased {
				x2 = ev.X
				y2 = ev.Y
			}

		case dt := <-ticker.C:
			screen.Debug(dt, 0, 0)
			t += 0.1
			screen.Clear()
			screen.DrawRect(int(math.Cos(t)*10)+50, 10, 5, 5, true, '#')
			screen.DrawLine(x1, y1, x2, y2, '@')
			screen.Display()
		}
	}
}
