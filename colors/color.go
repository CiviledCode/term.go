package colors

import "fmt"

// Color represents the various color formats that ASCII escape codes use for text-colorization.
type Color interface {
	// Compress condenses the color down into a string format that we can append to the end of an escape character.
	Compress() string
}

// DefaultColor represents the builtin constant colors.
type DefaultColor uint8

// Compress ...
func (d DefaultColor) Compress() string {
	return fmt.Sprintf(";%v", d)
}

const (
	// default colors.
	Black   DefaultColor = 0
	Red     DefaultColor = 1
	Green   DefaultColor = 2
	Yellow  DefaultColor = 3
	Blue    DefaultColor = 4
	Magenta DefaultColor = 5
	Cyan    DefaultColor = 6
	White   DefaultColor = 7
)
