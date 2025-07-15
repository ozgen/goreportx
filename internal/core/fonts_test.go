package core

import (
	"testing"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/assert"
)

type fontStyleCase struct {
	italic bool
	bold   bool
}

func (c fontStyleCase) TestName() string {
	switch {
	case c.bold && c.italic:
		return "BoldItalic"
	case c.bold:
		return "Bold"
	case c.italic:
		return "Italic"
	default:
		return "Regular"
	}
}

func TestFindFontPaths_ReturnsCorrectPaths_WhenFontsExist(t *testing.T) {
	// Act
	paths, err := findFontPaths()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, paths.Regular)
	assert.NotEmpty(t, paths.Bold)
	assert.NotEmpty(t, paths.Italic)
}

func TestSetFont_DoesNotPanic_WithAllStyles(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Define test cases
	cases := []fontStyleCase{
		{false, false},
		{true, false},
		{false, true},
		{true, true},
	}

	for _, c := range cases {
		t.Run(c.TestName(), func(t *testing.T) {
			SetFont(pdf, c.italic, c.bold)
		})
	}
}
