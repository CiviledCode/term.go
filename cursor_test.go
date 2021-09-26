package go_colorize

import (
	"github.com/civiledcode/term.go/screen"
	"testing"
)

func TestSetPosition(t *testing.T) {
	screen := screen.NewScreen()
	x, y := screen.Size()
	screen.ScreenCursor.SetPosition(uint16(x/2), uint16(y/2))
	screen.Close()
}
