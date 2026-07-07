// banner.go generates dot-matrix ASCII banners for the terminal UI.
//
// The dot-matrix style is a signature visual element of Moodwave CLI,
// defined in cli_design.md as: "A dot-matrix look for headers, mood labels,
// track names, status chips, and empty-state messages."
//
// Each character is rendered as a 5×5 grid of dots, scalable to any width.
package visuals

import "strings"

// dotMatrix maps each character to a 5×5 bit grid (top row first).
// 1 = filled, 0 = space. Each row is a 5-bit value.
var dotMatrix = map[rune][5]uint8{
	'A': {0b01110, 0b10001, 0b11111, 0b10001, 0b10001},
	'B': {0b11110, 0b10001, 0b11110, 0b10001, 0b11110},
	'C': {0b01111, 0b10000, 0b10000, 0b10000, 0b01111},
	'D': {0b11110, 0b10001, 0b10001, 0b10001, 0b11110},
	'E': {0b11111, 0b10000, 0b11110, 0b10000, 0b11111},
	'F': {0b11111, 0b10000, 0b11110, 0b10000, 0b10000},
	'G': {0b01111, 0b10000, 0b10011, 0b10001, 0b01111},
	'H': {0b10001, 0b10001, 0b11111, 0b10001, 0b10001},
	'I': {0b11111, 0b00100, 0b00100, 0b00100, 0b11111},
	'J': {0b11111, 0b00010, 0b00010, 0b10010, 0b01100},
	'K': {0b10001, 0b10010, 0b11100, 0b10010, 0b10001},
	'L': {0b10000, 0b10000, 0b10000, 0b10000, 0b11111},
	'M': {0b10001, 0b11011, 0b10101, 0b10001, 0b10001},
	'N': {0b10001, 0b11001, 0b10101, 0b10011, 0b10001},
	'O': {0b01110, 0b10001, 0b10001, 0b10001, 0b01110},
	'P': {0b11110, 0b10001, 0b11110, 0b10000, 0b10000},
	'Q': {0b01110, 0b10001, 0b10101, 0b10010, 0b01101},
	'R': {0b11110, 0b10001, 0b11110, 0b10010, 0b10001},
	'S': {0b01111, 0b10000, 0b01110, 0b00001, 0b11110},
	'T': {0b11111, 0b00100, 0b00100, 0b00100, 0b00100},
	'U': {0b10001, 0b10001, 0b10001, 0b10001, 0b01110},
	'V': {0b10001, 0b10001, 0b10001, 0b01010, 0b00100},
	'W': {0b10001, 0b10001, 0b10101, 0b11011, 0b10001},
	'X': {0b10001, 0b01010, 0b00100, 0b01010, 0b10001},
	'Y': {0b10001, 0b01010, 0b00100, 0b00100, 0b00100},
	'Z': {0b11111, 0b00010, 0b00100, 0b01000, 0b11111},
	'0': {0b01110, 0b10011, 0b10101, 0b11001, 0b01110},
	'1': {0b00100, 0b01100, 0b00100, 0b00100, 0b01110},
	'2': {0b01110, 0b10001, 0b00110, 0b01000, 0b11111},
	'3': {0b11111, 0b00010, 0b00110, 0b00001, 0b11110},
	'4': {0b00011, 0b00101, 0b01001, 0b11111, 0b00001},
	'5': {0b11111, 0b10000, 0b11110, 0b00001, 0b11110},
	'6': {0b01110, 0b10000, 0b11110, 0b10001, 0b01110},
	'7': {0b11111, 0b00001, 0b00010, 0b00100, 0b00100},
	'8': {0b01110, 0b10001, 0b01110, 0b10001, 0b01110},
	'9': {0b01110, 0b10001, 0b01111, 0b00001, 0b01110},
	'-': {0b00000, 0b00000, 0b11111, 0b00000, 0b00000},
	' ': {0b00000, 0b00000, 0b00000, 0b00000, 0b00000},
	'!': {0b00100, 0b00100, 0b00100, 0b00000, 0b00100},
	'.': {0b00000, 0b00000, 0b00000, 0b00000, 0b00100},
	'%': {0b11001, 0b11010, 0b00100, 0b01011, 0b10011},
	'~': {0b00000, 0b01001, 0b10110, 0b00000, 0b00000},
}

// renderDotMatrixBanner converts a string into a 5-row dot-matrix banner.
// The banner is returned as a slice of 5 strings, one per row.
// If noUnicode is true, filled dots use '#', otherwise '█'.
func renderDotMatrixBanner(text string, noUnicode bool) []string {
	upper := strings.ToUpper(text)
	runes := []rune(upper)

	filled := '█'
	empty := ' '
	if noUnicode {
		filled = '#'
	}

	rows := make([]strings.Builder, 5)

	for i, ch := range runes {
		if i > 0 {
			// Character spacing.
			for row := range rows {
				rows[row].WriteRune(empty)
				rows[row].WriteRune(empty)
			}
		}

		grid, ok := dotMatrix[ch]
		if !ok {
			// Unknown character — render as a space.
			for row := range rows {
				rows[row].WriteString("     ")
			}
			continue
		}

		for row := 0; row < 5; row++ {
			bits := grid[row]
			for bit := 4; bit >= 0; bit-- {
				if (bits>>uint(bit))&1 == 1 {
					rows[row].WriteRune(filled)
				} else {
					rows[row].WriteRune(empty)
				}
			}
		}
	}

	result := make([]string, 5)
	for i, row := range rows {
		result[i] = row.String()
	}
	return result
}
