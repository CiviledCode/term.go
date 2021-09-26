package go_colorize

import (
	"fmt"
	"github.com/civiledcode/term.go/screen"
	"testing"
)

func TestScreenWrite(t *testing.T) {
	screen := screen.NewScreen()
	screen.ClearScreen()
	fmt.Fprintln(screen.Terminal, "Hello World!")
	screen.Close()
}

func TestScreenClearing(t *testing.T) {
	screen := screen.NewScreen()
	screen.ScreenCursor.SetPosition(10, 10)
	screen.ClearToCursor(true)
	screen.Close()
}
