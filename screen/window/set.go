package window

import _ "encoding/json"

const (
	Pixel SetUnit = iota
	Percentage
)

// SetUnit represents a unit of measurement that a set can quantify.
type SetUnit uint8

// Set represents a group of numbers quantified in different ways.
type Set struct {
	// X represents the X quantity in 'Unit' units.
	X uint16 `json:"x"`

	// Y represents the Y quantity in 'Unit' units.
	Y uint16 `json:"y"`

	// Unit represents the unit of measurement for the set. This is used internally within Value to calculate the value.
	Unit SetUnit `json:"unit"`
}

// Value produces the raw pixel value for the set. This requires the width and height of the screen for scalar quantities.
func (s Set) Value(width, height uint16) (uint16, uint16) {
	switch s.Unit {
	case Pixel:
		return s.X, s.Y
	case Percentage:
		widthPerPercentage := float64(width) / 100
		heightPerPercentage := float64(height) / 100
		return uint16(widthPerPercentage * float64(s.X)), uint16(heightPerPercentage * float64(s.Y))
	}

	return 0, 0
}
