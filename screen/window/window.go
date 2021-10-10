package window

import (
	_ "encoding/json"
	"fmt"
	"github.com/civiledcode/term.go/colors"
	"github.com/civiledcode/term.go/screen"
	"log"
	"strings"
)

// Window represents a rectangular zone within the terminal. Windows are typically responsible for drawing
// content to said areas using events.
type Window struct {
	// Size represents the width and height of the window.
	Size Set `json:"size"`

	// Position represents the X and Y coordinate of the top left corner of the window.
	Position Set `json:"position"`

	// TODO: Remove this field
	// Color represents the background color of the blank buffer space being drawn.
	Color colors.Color

	// Name is the name of the window.
	Name string `json:"name"`

	// ID is the unique ID of the window. This should be a human-readable format.
	ID string `json:"id"`

	// FocusUpdated depicts if a window needs to be focused for an update call on its elements to be called.
	FocusUpdated bool `json:"focus_updated"`

	// FocusPriority depicts the order in which this window is focused. The larger the number, the least
	// in priority when cycling through windows. If this number is a negative, we don't ever try to focus on this window.
	FocusPriority int8 `json:"focus_priority"`

	// Visible depicts if the window should be drawn or not.
	Visible bool `json:"visible"`

	// events represents a list of events mapped to callback functions for Emitting events.
	events map[Event]map[uint8]EventListener

	// focused depicts if the window is currently in focus.
	focused bool

	testString string
}

// Focus is called when focus is switched to this window.
func (w *Window) Focus(s *screen.Screen) {
	w.focused = true
	s.ScreenCursor.SetPosition(w.Position.Value(s.Size()))
}

// Unfocus is called when focus is shifted away from this window.
func (w *Window) Unfocus(_ *screen.Screen) {
	w.focused = false
}

// Emit executes all functions registered under the ID defined in event.
// The params are arguments passed to the event.
func (w *Window) Emit(event Event, params ...interface{}) {
	log.Printf("Event Called in '%v': Event:%v Params:%v\n", w.Name, event, params)
	events := w.events[event]

	if len(events) == 0 {
		return
	}

	go func() {
		if len(events) > 1 {
			// TODO: Sort the events by their priority and execute them in this order.

			for _, e := range events {
				e(event, params)
			}
		} else {
			for _, e := range events {
				e(event, params)
			}
		}
	}()
}

// OnEvent binds a callback function to events.
func (w *Window) OnEvent(listener EventListener, priority uint8, events ...Event) {
	for _, event := range events {
		w.events[event][priority] = listener
	}
}

// RemoveEvent removes a listener from the event list using the priority and event name.
func (w *Window) RemoveEvent(event Event, priority uint8) {
	delete(w.events[event], priority)
}

// update is called everytime a key is pressed and receives the key pressed.
// This generally should handle logic regarding key presses within windows.
func (w *Window) update(screen *screen.Screen, input rune, isCtrl bool) {
	w.Emit(UpdateEvent, input, isCtrl)
	if !isCtrl {
		w.testString = fmt.Sprintf("%c", input)
	}
}

// draw is called everytime a key is pressed.
// This generally should handle rendering the window to the screen.
func (w *Window) draw(screen *screen.Screen) {
	w.Emit(DrawEvent)
	wWidth, wHeight := screen.Size()
	wHeight--
	width, height := w.Size.Value(wWidth, wHeight)
	screen.ScreenCursor.Save()
	x, y := w.Position.Value(wWidth, wHeight)

	if y == 0 {
		y++
	}

	bodyStyle := colors.NewStyle(colors.Normal)
	bodyStyle.SetBackgroundColor(w.Color)

	for i := uint16(0); i < height; i++ {
		screen.ScreenCursor.SetPosition(x, y+i)
		// Render the title
		if i == 0 {
			screen.ResetStyle()
			titleStyle := colors.NewStyle(colors.Bold)

			if w.focused {
				titleStyle.SetForegroundColor(colors.Green)
			} else {
				titleStyle.SetForegroundColor(colors.Red)
			}

			_, err := fmt.Fprintf(screen.Terminal, "%v\n", titleStyle.Colorize(centerText(w.Name, int(width))))

			if err != nil {
				log.Fatalln(err)
			}
		} else {
			_, err := fmt.Fprintf(screen.Terminal, "%v\n", bodyStyle.Colorize(strings.Repeat(w.testString, int(width))))

			if err != nil {
				log.Fatalln(err)
			}
		}
		screen.ScreenCursor.Return()
	}
}

// Focused depicts if the window was focused or not.
func (w *Window) Focused() bool {
	return w.focused
}

// centerText adds spaces surrounding text to fulfill len(text) == width.
func centerText(text string, width int) string {
	spacesNeeded := width - len(text)

	if spacesNeeded%2 == 1 {
		// Odd number
		spaces := strings.Repeat(" ", spacesNeeded/2)
		return spaces + text + spaces + " "
	} else if spacesNeeded <= 0 {
		return text
	} else {
		// Even number. Simple maths.
		spaces := strings.Repeat(" ", spacesNeeded/2)
		return spaces + text + spaces
	}
}
