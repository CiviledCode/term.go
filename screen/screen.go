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

func (s *Screen) JumpLines() {
	fmt.Fprint(s.Terminal, "\x1b[2J")
}

// ResetStyle wipes the current color configuration for the next characters being drawn.
func (s *Screen) ResetStyle() {
	fmt.Fprint(s.Terminal, "\x1b[0m")
}

// Size retrieves the size of the working Terminal space.
func (s *Screen) Size() (int, int) {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
		return -1, -1
	}
	return int(ws.Col), int(ws.Row)
}

// ClearScreenHistory sends an ASCII escape character to the Screen writer to clear the terminal history.
func (s *Screen) ClearScreenHistory() {
	fmt.Fprint(s.Terminal, "\\\033c")
}

// ClearLines clears all the visible lines on the screen.
func (s *Screen) ClearLines() {
	fmt.Fprint(s.Terminal, "\u001b[2J")
}

// Up moves the screen up x amount of lines.
func (s *Screen) Up(amount uint8) {
	fmt.Fprintf(s.Terminal, "\x1b[%vA", amount)
}

// Down moves the screen down x amount of lines.
func (s *Screen) Down(amount uint8) {
	fmt.Fprintf(s.Terminal, "\x1b[%vB", amount)
}

// ClearLine sends an ASCII escape character to the Screen writer to clear the current line that the cursor is on.
func (s *Screen) ClearLine() {
	fmt.Fprint(s.Terminal, "\u001b[2K")
}

// ClearFromCursor sends an ASCII escape character to the Screen writer
func (s *Screen) ClearFromCursor(eof bool) {
	if !eof {
		fmt.Fprint(s.Terminal, "\u001B[0K")
	} else {
		fmt.Fprint(s.Terminal, "\u001B[0J")
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
	s.ClearLines()
	s.ScreenCursor.Home()

	err := s.Terminal.Restore()
	if err != nil {
		panic(err)
	}
}

// ShouldClose depicts if the application has had a close call or not.
func (s *Screen) ShouldClose() bool {
	return !s.running
}
