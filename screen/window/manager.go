package window

import (
	"fmt"
	"github.com/civiledcode/term.go/screen"
	"github.com/pkg/term"
	"io"
	"log"
	"sort"
	"strings"
)

var (
	STDERR io.Writer

	// QuitKey represents the ctrl combination key to quit out of the program.
	QuitKey rune = 81
	// NextWindow represents the ctrl combination key to move to the next window in the window priority list.
	NextWindow rune = 91
	// PrevWindow represents the ctrl combination key to move to the previous window in the window priority list.
	PrevWindow rune = 93
	// UnfocusWindow represents the ctrl combination key to deselect the current window (no window focus).
	UnfocusWindow rune = 72
)

const debugConsoleLocation = "/dev/pts/1"

// TODO: Window event manager to add window drawing events and pass in a buffer.
// Doing things this way will allow people to still modularize their code, but will decrease the need of
// object oriented programming and go down a more data oriented route to make layouts easier later.

// EVENTS: OnDraw(buffer), OnUpdate(input, isCtrl)

// Add an event listener and an event handler. The handler is to handle all events for a specific window.
// The listener is meant to handle all subscribed events for a specific window. (Reference syscall lib)

// TODO: Add a focus lock channel bool that's passed as a param of focus to depict if we want
// to lock the focus iteration functions to allow dialogue controls and panic windows.

func init() {
	// TODO: Remove the debug console later.
	t, err := term.Open(debugConsoleLocation, term.Speed(19200))
	if err != nil {
		fmt.Println(err)
	}

	STDERR = t

	log.SetOutput(STDERR)
}

// windowManager is responsible for adding and updating console elements.
type windowManager struct {
	console *screen.Screen

	windows map[int8]*Window
}

// NewWindowManager creates a new windowManager using a terminal screen.
func NewWindowManager(s *screen.Screen) *windowManager {
	m := &windowManager{console: s, windows: make(map[int8]*Window)}
	go handleResizes(m)
	return m
}

// Update should be called every tick in order to update all child objects within the screen.
func (m *windowManager) Update(content rune, isCtrl bool) {
	m.console.ClearLines()

	// TODO: Handle this code better (Switch statement?)
	// TODO: Unhardcode keybinds.

	// Handle input
	if isCtrl {
		if content == QuitKey {
			m.console.Close()
			m.console.ClearScreenHistory()
			return
		} else if content == NextWindow || content == PrevWindow {
			// Cycle Windows.
			// CTRL + [ or ]
			focused, index, l := m.Focused()

			// If no window is focused, we should default to the first window.
			if focused == nil {
				// TODO: Create a mode to toggle
				focused = m.windows[int8(l[0])]
				focused.Focus(m.console)
			} else {
				// If a window is focused, we should determine the direction and increment to the next window.
				focused.Unfocus(m.console)

				if content == NextWindow {
					// Depict if we should start at the least prioritized or not.
					if index == 0 {
						w := m.windows[int8(l[len(l)-1])]
						w.Focus(m.console)
					} else {
						w := m.windows[int8(l[index-1])]
						w.Focus(m.console)
					}
				} else {
					// Depict if we should start at 0 or not.
					if index == len(l)-1 {
						w := m.windows[int8(l[0])]
						w.Focus(m.console)
					} else {
						w := m.windows[int8(l[index+1])]
						w.Focus(m.console)
					}
				}
			}
		} else if content == UnfocusWindow {
			// Unfocus from current window.
			// CTRL + Backspace
			focused, _, _ := m.Focused()
			if focused != nil {
				focused.Unfocus(m.console)
			}
		}
	}
	// Window calls.
	for _, w := range m.windows {
		if w.focused || !w.FocusUpdated {
			w.update(m.console, content, isCtrl)
		}

		if w.Visible {
			w.draw(m.console)
		}
	}
}

// GetWindow retrieves a window and its focus priority.
func (m *windowManager) GetWindow(id string) *Window {
	for _, window := range m.windows {
		if strings.EqualFold(window.ID, id) {
			return window
		}
	}

	return nil
}

// AddWindow adds a new child window to the windowManager to be updated.
func (m *windowManager) AddWindow(w *Window) {
	if m.windows[w.FocusPriority] != nil {
		m.console.Close()
		panic(fmt.Sprintf("More than one window attempted to register under focus index %v.", w.FocusPriority))
	} else {
		m.windows[w.FocusPriority] = w
	}
}

// Focused retrieves the focused window, it's index within the simplified focus list, and the simplified focus list.
// If no window is focused, we return nil and -1. Windows with a negative focus index aren't returned, as they
// should never be in focus.
func (m *windowManager) Focused() (*Window, int, []int) {
	var focused *Window
	var focusedIndex int
	l := make([]int, 0, len(m.windows))

	// Load all the values into a list and find the focused window.
	for key, w := range m.windows {
		if key < 0 {
			continue
		}

		if w.Focused() {
			focused = w
			focusedIndex = int(key)
		}
		l = append(l, int(key))
	}

	// Sort the list.
	sort.Ints(l)

	for i := 0; i < len(l); i++ {
		if l[i] == focusedIndex {
			return focused, i, l
		}
	}

	return nil, -1, l
}

func handleResizes(m *windowManager) {
	for {
		m.console.InputManager.Resize()
		m.Update(19132, true)
	}
}
