package main

import (
	"math"
	"time"

	"github.com/mslaursen/terminux"
)

type Vec3 struct {
	X, Y, Z float64
}

type Vec2 struct {
	X, Y int
}

func main() {
	screen := terminux.NewScreenDefault()
	screen.HideCursor()
	defer screen.Restore()

	go screen.ListenForEvents()

	ticker := screen.Ticker(time.Second / 60)
	defer ticker.Stop()

	w, h := screen.Size()

	cube := []Vec3{
		{-1, -1, -1},
		{1, -1, -1},
		{1, 1, -1},
		{-1, 1, -1},
		{-1, -1, 1},
		{1, -1, 1},
		{1, 1, 1},
		{-1, 1, 1},
	}

	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 0},
		{4, 5}, {5, 6}, {6, 7}, {7, 4},
		{0, 4}, {1, 5}, {2, 6}, {3, 7},
	}

	angle := 0.0

	c := screen.ResizeChan()

	for {
		select {
		case ev := <-screen.Events():
			if ev.Type == terminux.KeyPressed && ev.Key == "q" {
				return
			}
			screen.Debug(ev.Type, 0, 2)
		case <-c:
			w, h = screen.Resize()
		case <-ticker.C:
			angle += 0.03

			screen.Clear()

			projected := make([]Vec2, len(cube))

			for i, v := range cube {
				ry := rotateY(v, angle)
				rx := rotateX(ry, angle*0.7)

				dist := 3.0
				scale := 40.0
				z := rx.Z + dist

				px := (rx.X / z) * scale
				py := (rx.Y / z) * scale

				projected[i] = Vec2{
					X: int(px) + w/2,
					Y: int(py) + h/2,
				}
			}

			for _, e := range edges {
				a := projected[e[0]]
				b := projected[e[1]]

				screen.DrawLine(
					a.X+int(math.Cos(angle)*40), a.Y,
					b.X+int(math.Cos(angle)*40), b.Y,
					terminux.PixelFull,
					terminux.Cyan,
					terminux.Cyan,
				)
			}

			screen.Display()
		}
	}
}

func rotateX(v Vec3, angle float64) Vec3 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Vec3{
		X: v.X,
		Y: v.Y*c - v.Z*s,
		Z: v.Y*s + v.Z*c,
	}
}

func rotateY(v Vec3, angle float64) Vec3 {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Vec3{
		X: v.X*c + v.Z*s,
		Y: v.Y,
		Z: -v.X*s + v.Z*c,
	}
}
