// File: pkg/core/factory.go
package core

// RendererFactory encapsulates all the configuration options needed
// to construct a reusable PDF Renderer instance. It serves as a flexible
// builder for generating PDF documents with optional background, header,
// and footer images.
//
// This abstraction helps separate layout styling and rendering behavior
// from the data model and rendering logic.
type RendererFactory struct {
	// FontSizes defines the font size configuration for various text elements (e.g., H1, H2, body).
	FontSizes FontSizes

	// ShowPageNumber indicates whether page numbers should be displayed on each page.
	ShowPageNumber bool

	// Base64Background is a base64-encoded string representing the background image (optional).
	Base64Background string

	// Base64Header is a base64-encoded string representing the header image (optional).
	Base64Header string

	// Base64Footer is a base64-encoded string representing the footer image (optional).
	Base64Footer string
}

// NewRendererFactory creates a RendererFactory instance with default font sizes
// and page number enabled.
func NewRendererFactory() *RendererFactory {
	return &RendererFactory{
		FontSizes: FontSizes{
			H1:     24,
			H2:     18,
			H3:     14,
			P:      12,
			Footer: 10,
		},
		ShowPageNumber: true,
	}
}

// WithFontSizes sets custom font sizes for headings, paragraph, and footer.
func (f *RendererFactory) WithFontSizes(sizes FontSizes) *RendererFactory {
	f.FontSizes = sizes
	return f
}

// WithPageNumbers toggles whether page numbers should appear on rendered PDFs.
func (f *RendererFactory) WithPageNumbers(show bool) *RendererFactory {
	f.ShowPageNumber = show
	return f
}

// WithBaseImage sets a base64-encoded background image (optional).
func (f *RendererFactory) WithBaseImage(base64 string) *RendererFactory {
	f.Base64Background = base64
	return f
}

// WithHeaderImage sets a base64-encoded header image (optional).
func (f *RendererFactory) WithHeaderImage(base64 string) *RendererFactory {
	f.Base64Header = base64
	return f
}

// WithFooterImage sets a base64-encoded footer image (optional).
func (f *RendererFactory) WithFooterImage(base64 string) *RendererFactory {
	f.Base64Footer = base64
	return f
}

// Build creates a new Renderer instance based on the current configuration.
// If no images are provided, it defaults to a simple text-based renderer.
// Otherwise, it returns a renderer with the specified base64-encoded images.
func (f *RendererFactory) Build() (*Renderer, error) {
	if f.Base64Background == "" && f.Base64Header == "" && f.Base64Footer == "" {
		return NewRenderer(f.FontSizes, f.ShowPageNumber)
	}
	return NewRendererWithBase64Images(
		f.Base64Background,
		f.Base64Header,
		f.Base64Footer,
		f.FontSizes,
		f.ShowPageNumber,
	)
}
