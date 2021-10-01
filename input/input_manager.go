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

	resizeHandler func(uint16, uint16) error

	terminal *term.Term
}

// NewManager creates a new input manager and subscribes to its signals.
func NewManager(t *term.Term) *Manager {
	im := Manager{signalChannel: make(chan os.Signal, 1), terminal: t}
	signal.Notify(im.signalChannel, unix.SIGWINCH)
	return &im
}

// OnResize sets the event listener that's executed when the window is resized.
func (im *Manager) OnResize(event func(uint16, uint16) error) {
	im.resizeHandler = event
}

// Update should be called every tick or frame of your application to get the latest system calls and inputs.
func (im *Manager) Update() {
	// TODO: Unhold the thread
	s := <-im.signalChannel
	switch s {
	case unix.SIGWINCH:
		{
			if im.resizeHandler != nil {
				ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
				if err != nil {
					panic(err)
				}

				err = im.resizeHandler(ws.Col, ws.Row)

				if err != nil {
					panic(err)
				}
			}
		}
	default:
		fmt.Println("Unknown signal.")
	}
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
