package core

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

// FontPaths represents the file paths for the regular, bold, and italic font variants.
type FontPaths struct {
	Regular string
	Bold    string
	Italic  string
}
