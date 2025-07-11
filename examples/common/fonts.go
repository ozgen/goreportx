package common

import "github.com/ozgen/goreportx/internal/renderer"

// DefaultFontSizes provides standard font sizes for typical reports.
var DefaultFontSizes = renderer.FontSizes{
	H1:     24,
	H2:     20,
	H3:     16,
	P:      12,
	Footer: 10,
}

// LargeFontSizes provides larger font sizes for more visual emphasis.
var LargeFontSizes = renderer.FontSizes{
	H1:     30,
	H2:     24,
	H3:     18,
	P:      14,
	Footer: 12,
}
