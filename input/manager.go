package input

import (
	"fmt"
	"github.com/pkg/term"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

// Manager interacts with signals to provide a manager for all user input within the terminal.
type Manager struct {
	signalChannel chan os.Signal

	terminal *term.Term
}

// NewManager creates a new input manager and subscribes to its signals.
func NewManager(t *term.Term) *Manager {
	im := Manager{signalChannel: make(chan os.Signal, 1), terminal: t}
	signal.Notify(im.signalChannel, unix.SIGWINCH)
	return &im
}

// Resize holds the thread until a resize event happens. If a resize event happens, the new window sizes are returned.
func (im *Manager) Resize() (uint16, uint16) {
	s := <-im.signalChannel
	switch s {
	case unix.SIGWINCH:
		{
			ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
			if err != nil {
				panic(err)
			}

			return ws.Col, ws.Row
		}
	default:
		fmt.Println("Unknown signal.")
	}

	return 0, 0
}

// Input returns the latest input being pressed. This will hold the current thread until a new key is pressed,
// so using it as means of limiting update calls has proven effective.
func (im *Manager) Input() (rune, bool) {
	l := make([]byte, 8)
	_, err := im.terminal.Read(l)

	if err != nil {
		fmt.Println(err)
	}

	return parseKey(l)
}
