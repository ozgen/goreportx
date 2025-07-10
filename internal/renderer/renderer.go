// File: renderer/renderer.go
package renderer

import (
	"github.com/signintech/gopdf"
)

const (
	pageHeight   = 842.0
	pageMargin   = 50.0
	footerHeight = 30.0
	contentLimit = pageHeight - pageMargin - footerHeight
)

type TextChunk struct {
	Text   string
	Italic bool
	Bold   bool
}

type FontSizes struct {
	H1     float64
	H2     float64
	H3     float64
	P      float64
	Footer float64
}

type Renderer struct {
	pdf            *gopdf.GoPdf
	y              float64
	pageWidth      float64
	footerText     string
	pageNumber     int
	showPageNumber bool
	backgroundImg  string
	headerImg      string
	footerImg      string
	FontSize       FontSizes
}
