package terminux

const (
	// ansiEsc             = "\x1b"
	ansiCursorMove      = "\033[%d;%dH"
	ansiHideCursor      = "\033[?25l"
	ansiShowCursor      = "\033[?25h"
	ansiClearScreen     = "\033[2J"
	ansiCursorHome      = "\033[H"
	ansiEnableMouse     = "\x1b[?1000h"
	ansiDisableMouse    = "\x1b[?1000l"
	ansiEnableMouseSGR  = "\x1b[?1006h"
	ansiDisableMouseSGR = "\x1b[?1006l"
)

const (
	Black   = "\u001b[30m"
	Red     = "\u001b[31m"
	Green   = "\u001b[32m"
	Yellow  = "\u001b[33m"
	Blue    = "\u001b[34m"
	Magenta = "\u001b[35m"
	Cyan    = "\u001b[36m"
	White   = "\u001b[37m"

	BrightBlack   = "\u001b[90m"
	BrightRed     = "\u001b[91m"
	BrightGreen   = "\u001b[92m"
	BrightYellow  = "\u001b[93m"
	BrightBlue    = "\u001b[94m"
	BrightMagenta = "\u001b[95m"
	BrightCyan    = "\u001b[96m"
	BrightWhite   = "\u001b[97m"

	BgBlack   = "\u001b[40m"
	BgRed     = "\u001b[41m"
	BgGreen   = "\u001b[42m"
	BgYellow  = "\u001b[43m"
	BgBlue    = "\u001b[44m"
	BgMagenta = "\u001b[45m"
	BgCyan    = "\u001b[46m"
	BgWhite   = "\u001b[47m"

	BgBrightBlack   = "\u001b[100m"
	BgBrightRed     = "\u001b[101m"
	BgBrightGreen   = "\u001b[102m"
	BgBrightYellow  = "\u001b[103m"
	BgBrightBlue    = "\u001b[104m"
	BgBrightMagenta = "\u001b[105m"
	BgBrightCyan    = "\u001b[106m"
	BgBrightWhite   = "\u001b[107m"

	Reset = "\u001b[0m"
)

const (
	PixelFull      = '█'
	PixelDark      = '▓'
	PixelMed       = '▒'
	PixelLight     = '░'
	PixelUpperHalf = '▀'
	PixelLowerHalf = '▄'
)
