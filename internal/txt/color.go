package txt

type Color string

var ColorReset Color = "\x1b[0m"
var ColorBgBlack Color = "\x1b[40m"
var Colors = []Color{
	"\x1b[41m", // Red
	"\x1b[42m", // Green
	"\x1b[43m", // Yellow
	"\x1b[44m", // Blue
	"\x1b[45m", // Purple
	"\x1b[46m", // Cyan
	"\x1b[47m", // White
}
