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

	ticker := screen.Ticker(time.Second / 60)
	defer ticker.Stop()

	w, h := screen.Size()

	t := 0.1

	for {
		select {
		case ev := <-screen.Events():
			if ev.Type == terminux.KeyPressed && ev.Key == "q" {
				return
			}

		case <-ticker.C:
			t += 0.1
			x := math.Cos(t) * 10
			y := math.Sin(t) * 5
			screen.Clear()
			screen.Draw(int(x)+w/2, int(y)+h/2, terminux.PixelFull, terminux.Red, terminux.Red)
			screen.DrawRect(0, 0, 50, h, true, terminux.PixelFull, terminux.Red, terminux.Red)
			screen.DrawLine(0, 0, w, h, terminux.PixelDark, terminux.Blue, terminux.Blue)
			screen.Display()
		}
	}
}
