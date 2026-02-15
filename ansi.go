package terminux

const (
	// ansiEsc             = "\x1b"
	ansiCursorMove      string = "\033[%d;%dH"
	ansiHideCursor      string = "\033[?25l"
	ansiShowCursor      string = "\033[?25h"
	ansiClearScreen     string = "\033[2J"
	ansiCursorHome      string = "\033[H"
	ansiEnableMouse     string = "\x1b[?1000h"
	ansiDisableMouse    string = "\x1b[?1000l"
	ansiEnableMouseSGR  string = "\x1b[?1006h"
	ansiDisableMouseSGR string = "\x1b[?1006l"
)

const (
	Black   uint32 = 30
	Red     uint32 = 31
	Green   uint32 = 32
	Yellow  uint32 = 33
	Blue    uint32 = 34
	Magenta uint32 = 35
	Cyan    uint32 = 36
	White   uint32 = 37

	BrightBlack   uint32 = 90
	BrightRed     uint32 = 91
	BrightGreen   uint32 = 92
	BrightYellow  uint32 = 93
	BrightBlue    uint32 = 94
	BrightMagenta uint32 = 95
	BrightCyan    uint32 = 96
	BrightWhite   uint32 = 97

	BgBlack        uint32 = 40
	BgRed          uint32 = 41
	BgGreen        uint32 = 42
	BgYellow       uint32 = 43
	BgBlue         uint32 = 44
	BgMagentaColor uint32 = 45
	BgCyan         uint32 = 46
	BgWhite        uint32 = 47

	BgBrightBlack   uint32 = 100
	BgBrightRed     uint32 = 101
	BgBrightGreen   uint32 = 102
	BgBrightYellow  uint32 = 103
	BgBrightBlue    uint32 = 104
	BgBrightMagenta uint32 = 105
	BgBrightCyan    uint32 = 106
	BgBrightWhite   uint32 = 107

	Reset uint32 = 0
)

const (
	PixelFull      = '█'
	PixelDark      = '▓'
	PixelMed       = '▒'
	PixelLight     = '░'
	PixelUpperHalf = '▀'
	PixelLowerHalf = '▄'
)
