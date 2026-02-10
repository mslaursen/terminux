package main

import (
	"fmt"
	"math"
	"time"

	"github.com/mslaursen/terminux"
)

func main() {
	screen := terminux.NewScreenDefault()
	screen.HideCursor()
	defer screen.Restore()

	running := true

	screen.SetEventListener(func(event string) {
		if event == "q" {
			running = false
		}
		fmt.Println(event)

	})

	go screen.ListenForEvents()

	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	w, h := screen.GetSize()
	t := 0.0

	for range ticker.C {

		if !running {
			break
		}

		screen.Clear()

		t += 0.1

		x := math.Cos(t)*20 + (float64(w / 2))
		y := math.Sin(t)*8 + (float64(h / 2))

		screen.DrawRect(int(x), int(y), 5, 5, true, 'x')
		screen.DrawRect(int(x-5), int(y-5), 15, 15, false, '!')
		screen.Draw(int(x), int(y), '#')
		screen.DrawLine(20, 20, w, h, '-')
		screen.DrawRect(0, 0, w, h, false, '!')

		screen.Display()
	}
}
