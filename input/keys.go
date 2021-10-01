package input

import "unicode"

// parseKey receives a byte array of input sent by the terminal and converts it into a rune.
// The returned rune represents the display key pressed and the return bool represents if the key was a ctrl+ key or not.
// If we detect a function or special key, we return our own rune representing that key to allow us
// to hold all key press possibilities in a rune format. If the key is a ctrl+ key, then we return the letter key.
func parseKey(input []byte) (rune, bool) {
	if input[1] == 0 {
		k := rune(input[0])

		if unicode.IsGraphic(k) {
			// This key press is nothing fancy
			return k, false
		} else if unicode.IsControl(k) {

			// TODO: Improve accuracy
			return k + 64, true
		}
	} else {
		// Special keys.
		if input[0] == 27 {
			// TODO: Handle special keys
		} else {
			// TODO: HANDLE UNKNOWN KEY
		}
	}

	// This should never happen
	return 0, false
}
