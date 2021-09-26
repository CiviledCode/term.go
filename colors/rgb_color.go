package colors

import "fmt"

// RGBColor represents a custom 24-bit RGB color.
type RGBColor struct {
	R, G, B uint8
}

// Compress ...
func (r RGBColor) Compress() string {
	return fmt.Sprintf(";2;%v;%v;%v", r.R, r.G, r.B)
}
