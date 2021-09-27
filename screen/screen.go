package screen

import (
	"fmt"
	"github.com/civiledcode/term.go/input"
	"github.com/pkg/term"
	"golang.org/x/sys/unix"
	_ "unsafe"
)

// Screen represents a Terminal application Screen. This struct contains functions
// to manipulate and extend on top of the current Terminal using ascii escape codes.
type Screen struct {
	Terminal *term.Term

	running bool

	ScreenCursor Cursor

	InputManager *input.Manager
}

// NewScreen creates a new 'Screen' object using an io.Writer.
func NewScreen() *Screen {
	t, err := term.Open("/dev/tty", term.Speed(19200))
	if err != nil {
		fmt.Println(err)
	}

	err = term.RawMode(t)
	if err != nil {
		fmt.Println(err)
	}
	return &Screen{Terminal: t, ScreenCursor: Cursor{terminal: t}, running: true, InputManager: input.NewManager(t)}
}

// ClearScreen sends an ASCII escape character to the Screen writer to clear the text of all lines visible.
func (s *Screen) ClearScreen() {
	fmt.Fprint(s.Terminal, "\x1b[2J")
}

// Reset wipes the current color configuration for the next characters being drawn.
func (s *Screen) Reset() {
	fmt.Fprint(s.Terminal, "\x1b[0m")
}

// Size retrieves the size of the working Terminal space.
func (s *Screen) Size() (int, int) {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
		return -1, -1
	}
	return int(ws.Row), int(ws.Col)
}

// ClearLine sends an ASCII escape character to the Screen writer to clear the current line that the cursor is on.
func (s *Screen) ClearLine() {
	fmt.Fprint(s.Terminal, "\x1b[2K")
}

// ClearFromCursor sends an ASCII escape character to the Screen writer
func (s *Screen) ClearFromCursor(eof bool) {
	if !eof {
		fmt.Fprint(s.Terminal, "\x0b[1K")
	} else {
		fmt.Fprint(s.Terminal, "\x0b[1J")
	}
}

// ClearToCursor sends an ASCII escape character to the Screen writer to clear all content from a starting point to the cursor.
// If bos is set to true, we start at the beginning of the Screen and clear all content to the cursor. If not, we clear all content from the start of the line
// to the cursor.
func (s *Screen) ClearToCursor(bos bool) {
	if !bos {
		fmt.Fprint(s.Terminal, "\x1b[1K")
	} else {
		fmt.Fprint(s.Terminal, "\x1b[1J")
	}
}

// Close sets the Terminal to stop running.
func (s *Screen) Close() {
	s.running = false
	s.ClearScreen()

	err := s.Terminal.Restore()
	if err != nil {
		panic(err)
	}
}

// ShouldClose depicts if the application has had a close call or not.
func (s *Screen) ShouldClose() bool {
	return !s.running
}
