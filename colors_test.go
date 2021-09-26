package go_colorize

import (
	"github.com/civiledcode/term.go/colors"
	"testing"
)

func TestColors(t *testing.T) {
	normal := colors.NewStyle(colors.Normal)
	normal.SetForegroundColor(colors.Yellow)
	normal.SetBackgroundColor(colors.Black)
	t.Log(normal.Colorize("Hello World!"))

	boldStyle := colors.NewStyle(colors.Bold)
	boldStyle.SetForegroundColor(colors.Green)
	t.Log(boldStyle.Colorize("Hello World!"))

	dim := colors.NewStyle(colors.Dim)
	dim.SetBackgroundColor(colors.Red)
	t.Log(dim.Colorize("Hello World!"))

	underline := colors.NewStyle(colors.Underlined)
	underline.SetForegroundColor(colors.White)
	t.Log(underline.Colorize("Hello World!"))

	// this only works in some terminals
	flicker := colors.NewStyle(colors.Blink)
	flicker.SetBackgroundColor(colors.Magenta)
	t.Log(flicker.Colorize("Hello World!"))

	rgb := colors.NewStyle(colors.Bold)
	rgb.SetForegroundColor(colors.RGBColor{94, 224, 183})
	t.Log(rgb.Colorize("Hello World!"))
}
