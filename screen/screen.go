package screen

import (
	"fmt"
	"github.com/civiledcode/term.go/input"
	"github.com/pkg/term"
	"golang.org/x/sys/unix"
	_ "unsafe"
)

// screen represents a Terminal application screen. This struct contains functions
// to manipulate and extend on top of the current Terminal using ascii escape codes.
type screen struct {
	Terminal *term.Term

	running bool

	ScreenCursor Cursor

	InputManager *input.Manager
}

// NewScreen creates a new 'screen' object using an io.Writer.
func NewScreen() *screen {
	t, err := term.Open("/dev/tty", term.Speed(19200))
	if err != nil {
		fmt.Println(err)
	}

	err = term.RawMode(t)
	if err != nil {
		fmt.Println(err)
	}
	return &screen{Terminal: t, ScreenCursor: Cursor{terminal: t}, running: true, InputManager: input.NewManager(t)}
}

// ClearScreen sends an ASCII escape character to the screen writer to clear the text of all lines visible.
func (s *screen) ClearScreen() {
	fmt.Fprint(s.Terminal, "\x1b[2J")
}

// Reset wipes the current color configuration for the next characters being drawn.
func (s *screen) Reset() {
	fmt.Fprint(s.Terminal, "\x1b[0m")
}

// Size retrieves the size of the working Terminal space.
func (s *screen) Size() (int, int) {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
		return -1, -1
	}
	return int(ws.Col), int(ws.Row)
}

// ClearLine sends an ASCII escape character to the screen writer to clear the current line that the cursor is on.
func (s *screen) ClearLine() {
	fmt.Fprint(s.Terminal, "\x1b[2K")
}

// ClearFromCursor sends an ASCII escape character to the screen writer
func (s *screen) ClearFromCursor(eof bool) {
	if !eof {
		fmt.Fprint(s.Terminal, "\x0b[1K")
	} else {
		fmt.Fprint(s.Terminal, "\x0b[1J")
	}
}

// ClearToCursor sends an ASCII escape character to the screen writer to clear all content from a starting point to the cursor.
// If bos is set to true, we start at the beginning of the screen and clear all content to the cursor. If not, we clear all content from the start of the line
// to the cursor.
func (s *screen) ClearToCursor(bos bool) {
	if !bos {
		fmt.Fprint(s.Terminal, "\x1b[1K")
	} else {
		fmt.Fprint(s.Terminal, "\x1b[1J")
	}
}

// Close sets the Terminal to stop running.
func (s *screen) Close() {
	s.running = false
	s.ClearScreen()

	err := s.Terminal.Restore()
	if err != nil {
		panic(err)
	}
}

// ShouldClose depicts if the application has had a close call or not.
func (s *screen) ShouldClose() bool {
	return !s.running
}
