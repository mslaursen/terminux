
# terminux

`terminux` is a small Go library for drawing and handling input in a terminal using ANSI escape codes.

It gives you:

* Raw keyboard input
* Double buffering (only changed cells are redrawn)
* Simple drawing helpers (points, rectangles, lines, etc.)

No dependencies beyond `golang.org/x/term`.

---

## Installation

```bash
go get github.com/mslaursen/terminux
```

---

## Render Loop

1. Create a screen
2. Put the terminal into raw mode
3. Run a loop that:
   * clears the back buffer
   * draws things
   * flushes the buffer to the terminal
4. Restore the terminal on exit

Input is handled asynchronously via an event listener.

---

## Minimal example

```go
package main

import (
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
	})

	go screen.ListenForEvents()

	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()

	for range ticker.C {
		if !running {
			break
		}

		screen.Clear()
		screen.Draw(10, 10, '@')
		screen.Display()
	}
}
```


## Input

Input is read rune by rrune in raw mode.

```go
screen.SetEventListener(func(event string) {
	fmt.Println(event)
})
go screen.ListenForEvents()
```

Each key press is passed to the listener as a string.

---

## Terminal state

`terminux` switches the terminal into raw mode.
Always restore it before exiting:

```go
defer screen.Restore()
```

This clears the screen, shows the cursor again, and resets terminal settings.

---

