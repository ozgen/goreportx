// File: renderer/renderer.go
package pkg

import (
	"fmt"
	"github.com/signintech/gopdf"
)

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

// NewRenderer initializes a new PDF renderer with the specified font sizes and
// page number visibility. It loads standard Arial fonts and sets the default
// font size for body text.
//
// Returns a Renderer instance ready to write PDF content.
func NewRenderer(fontSizes FontSizes, showPageNumber bool) (*Renderer, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	fonts, err := findFontPaths()
	if err != nil {
		return nil, err
	}

	// Load fonts with styles
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Regular, gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		return nil, fmt.Errorf("regular font: %w", err)
	}
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Bold, gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		return nil, fmt.Errorf("bold font: %w", err)
	}
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Italic, gopdf.TtfOption{Style: gopdf.Italic}); err != nil {
		return nil, fmt.Errorf("italic font: %w", err)
	}
	if err := pdf.SetFont("Arial", "", fontSizes.P); err != nil {
		return nil, err
	}

	return &Renderer{
		pdf:            pdf,
		y:              50,
		pageNumber:     1,
		showPageNumber: showPageNumber,
		FontSize:       fontSizes,
	}, nil
}

// NewRendererWithBase64Images creates a new Renderer and overlays background,
// header, and footer images from base64 strings.
//
// The images are saved to temporary files and rendered immediately.
// This function is useful when generating branded or templated reports
// with header/footer banners.
//
// Parameters:
// - bgBase64: base64-encoded background image (optional)
// - headerBase64: base64-encoded header image (optional)
// - footerBase64: base64-encoded footer image (optional)
// - fontSizes: custom font size configuration
// - showPageNumber: whether to render page numbers
//
// Returns a fully initialized Renderer.
func NewRendererWithBase64Images(bgBase64, headerBase64, footerBase64 string, fontSizes FontSizes, showPageNumber bool) (*Renderer, error) {
	bgPath := mustSaveBase64ToTempFile(bgBase64)
	headerPath := mustSaveBase64ToTempFile(headerBase64)
	footerPath := mustSaveBase64ToTempFile(footerBase64)

	r, err := NewRenderer(fontSizes, showPageNumber)
	if err != nil {
		return nil, err
	}

	r.backgroundImg = bgPath
	r.headerImg = headerPath
	r.footerImg = footerPath

	// Draw background and header/footer images on the first page
	if bgPath != "" {
		_ = r.pdf.Image(bgPath, 0, 0, &gopdf.Rect{W: 595.28, H: 841.89})
	}
	if headerPath != "" {
		_ = r.pdf.Image(headerPath, 50, 20, &gopdf.Rect{W: 495.0, H: 40.0})
		r.y += 50
	}
	if footerPath != "" {
		_ = r.pdf.Image(footerPath, 50, 800, &gopdf.Rect{W: 495.0, H: 30.0})
	}
	return r, nil
}
