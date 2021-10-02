package screen

import (
	"fmt"
	"io"
)

// Cursor represents the position of a terminal cursor.
type Cursor struct {
	x, y uint16

	hidden bool

	terminal io.Writer
}

// Home moves the position of the cursor to 0,0.
func (c *Cursor) Home() {
	fmt.Fprintf(c.terminal, "\x1b[H")
}

// SetPosition sets the position of the cursor on a 2D coordinate plane. 0,0 is the top-left corner.
func (c *Cursor) SetPosition(x, y uint16) {
	c.x = x
	c.y = y
	fmt.Fprintf(c.terminal, "\x1b[%v;%vf", y, x)
}

// LinePosition sets the cursor position on the X axis.
func (c *Cursor) LinePosition(x uint16) {
	c.x = x
	fmt.Fprintf(c.terminal, "\x1b[%vG", x)
}

// Up moves the cursor up x amount of lines. If startOfLine is true, then the cursor is moved to the start of the resulting line.
func (c *Cursor) Up(amount uint16, startOfLine bool) {
	if startOfLine {
		c.SetPosition(uint16(0), c.y-amount)
	} else {
		c.SetPosition(c.x, c.y-amount)
	}

	c.y -= amount
}

// Down moves the cursor down x amount of lines. If startOfLine is true, then the cursor is moved to the start of the resulting line.
func (c *Cursor) Down(amount uint16, startOfLine bool) {
	if startOfLine {
		c.SetPosition(uint16(0), c.y+amount)
	} else {
		c.SetPosition(c.x, c.y+amount)
	}

	c.y += amount
}

// Right moves the cursor right x amount of spaces.
func (c *Cursor) Right(amount uint8) {
	fmt.Fprintf(c.terminal, "\x1b[%vC", amount)
	c.x += uint16(amount)
}

// Left moves the cursor left x amount of spaces.
func (c *Cursor) Left(amount uint8) {
	fmt.Fprintf(c.terminal, "\x1b[%vD", amount)
	c.x -= uint16(amount)
}

// Save saves the current position of the cursor. This value cannot yet be retrieved yet, so it's only useful for jumping to a section and jumping back quickly.
func (c *Cursor) Save() {
	fmt.Fprint(c.terminal, "\x1b[s")
}

// Return moves the position of the cursor to the last saved position.
func (c *Cursor) Return() {
	fmt.Fprint(c.terminal, "\x1b[u")
}

// Position returns the current position of the cursor represented on a 2D coordinate plane.
func (c *Cursor) Position() (uint16, uint16) {
	return c.x, c.y
}

// Show sets the cursor visibility to be shown.
func (c *Cursor) Show() {
	if c.hidden {
		c.hidden = false
		fmt.Fprint(c.terminal, "\x1b[?25h")
	}
}

// Hide sets the cursor visibility to be hidden.
func (c *Cursor) Hide() {
	if !c.hidden {
		c.hidden = true
		fmt.Fprint(c.terminal, "\x1b[?25l")
	}
}
