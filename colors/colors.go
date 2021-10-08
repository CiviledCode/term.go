package colors

import (
	"fmt"
)

// textMode is the mode that the text is going to be rendered in.
type textMode uint8

const (
	// text modes
	Normal     textMode = 0
	Bold       textMode = 1
	Dim        textMode = 2
	Underlined textMode = 4
	Blink      textMode = 5

	rgbForeground   textMode = 38
	rgbBackground   textMode = 48
	indexForeground textMode = 39
	indexBackground textMode = 49

	reset string = "\x1b[0m"
)

// TODO: Index colors

// Style creates and builds various formats of ASCII escape codes for text formatting in terminals.
type Style struct {
	mode textMode

	foreground Color
	background Color
}

// NewStyle creates a new style using the type of text you'd like to create.
func NewStyle(mode textMode) Style {
	return Style{mode: mode}
}

// SetForegroundColor sets the foreground color. Because all colors are parsed into strings, we can use a common interface to allow multiple color types.
func (s *Style) SetForegroundColor(c Color) {
	_, isDefault := c.(DefaultColor)
	if isDefault {
		if s.mode == rgbBackground || s.mode == indexBackground {
			s.background = nil
		}

		s.mode = Normal
		s.foreground = c
		return
	}

	rgb, isRGB := c.(RGBColor)
	if isRGB {
		s.mode = rgbForeground
		s.background = nil
		s.foreground = rgb
	}
}

// SetBackgroundColor sets the background color. Because all colors are parsed into strings, we can use a common interface to allow multiple color types.
func (s *Style) SetBackgroundColor(c Color) {
	_, isDefault := c.(DefaultColor)
	if isDefault {
		if s.mode == rgbForeground || s.mode == indexForeground {
			s.foreground = nil
		}

		s.mode = Normal
		s.background = c
		return
	}

	_, isRGB := c.(RGBColor)
	if isRGB {
		s.mode = rgbBackground
		s.background = c
	}
}

// Mode returns the mode of the style.
func (s *Style) Mode() textMode {
	return s.mode
}

// Colorize takes a string of text and stylizes it using the current style.
func (s *Style) Colorize(text string) string {
	c, mode := buildColor(s)
	if mode == indexForeground || mode == indexBackground {
		mode--
	}
	return fmt.Sprintf("\x1b[%v%vm%v%v", mode, c, text, reset)
}

// String ...
func (s Style) String() string {
	c, mode := buildColor(&s)
	if mode == indexForeground || mode == indexBackground {
		mode--
	}
	return fmt.Sprintf("\x1b[%v%vm", mode, c)
}

// buildColor converts multiple different color types into a single format that the console understands.
func buildColor(style *Style) (built string, mode textMode) {
	switch style.mode {
	case rgbForeground:
		{
			rgb := style.foreground.(RGBColor)

			built += rgb.Compress()
		}

	case rgbBackground:
		{
			rgb := style.background.(RGBColor)

			built += rgb.Compress()
		}

	default:
		{
			// Get the normal color
			if style.foreground != nil {
				c := style.foreground.(DefaultColor)
				c += 30
				built += c.Compress()
			}

			if style.background != nil {
				c := style.background.(DefaultColor)
				c += 40
				built += c.Compress()
			}
		}
	}
	return built, style.mode
}
