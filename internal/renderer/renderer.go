// Package renderer provides functionality for rendering structured HTML-like
// content into a PDF document using the gopdf library. It supports headers,
// footers, images, tables, styled text, and automatic pagination.
package renderer

import (
	"github.com/signintech/gopdf"
)

const (
	// pageHeight defines the total height of an A4 PDF page in points.
	pageHeight = 842.0

	// pageMargin defines the top and bottom margin of the page in points.
	pageMargin = 50.0

	// footerHeight defines the reserved height for the footer area.
	footerHeight = 30.0

	// contentLimit defines the maximum vertical space available for content on a single page.
	contentLimit = pageHeight - pageMargin - footerHeight
)

// TextChunk represents a piece of styled text with font style information.
type TextChunk struct {
	Text   string // The actual text content.
	Italic bool   // Whether the text is italicized.
	Bold   bool   // Whether the text is bold.
}

// FontSizes defines the font size configuration for various text elements in the PDF.
type FontSizes struct {
	H1     float64 // Font size for <h1> elements.
	H2     float64 // Font size for <h2> elements.
	H3     float64 // Font size for <h3> elements.
	P      float64 // Font size for paragraph and body text.
	Footer float64 // Font size for the footer section.
}

// Renderer is the main structure used to manage the layout and rendering of content into a PDF.
// It encapsulates the gopdf instance, font settings, current layout position, and additional options
// such as header/footer images and timestamp rendering.
type Renderer struct {
	pdf               *gopdf.GoPdf // Internal PDF instance from gopdf.
	y                 float64      // Current vertical position on the page.
	pageWidth         float64      // Width of the current page (default A4).
	footerText        string       // Footer text to be rendered on each page.
	pageNumber        int          // Current page number.
	showPageNumber    bool         // Whether to render the page number in the footer.
	backgroundImg     string       // Base64-encoded background image path.
	headerImg         string       // Base64-encoded header image path.
	footerImg         string       // Base64-encoded footer image path.
	FontSize          FontSizes    // Font size configuration for the document.
	TopRightTimestamp string       // Optional timestamp text to be shown at the top-right of each page.
}
