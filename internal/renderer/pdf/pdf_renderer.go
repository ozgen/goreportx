// File: renderer/pdf/pdf_renderer.go
package pdf

import (
	"bytes"
	"github.com/ozgen/goreportx/internal/interfaces"
	"github.com/ozgen/goreportx/internal/renderer"
	"html/template"
	"log"
	"os"
	"time"
)

// PDFRenderer renders a given report using an HTML template into a PDF file.
// It supports optional image-based headers, footers, background, and timestamp injection.
type PDFRenderer struct {
	Report            interface{}        // The report data passed to the template
	Template          *template.Template // The parsed HTML template for the PDF layout
	FontSizes         renderer.FontSizes // Font sizes to use in the PDF rendering
	UseImages         bool               // If true, use header/footer/base images
	HeaderImgBase64   string             // Base64-encoded header image
	FooterImgBase64   string             // Base64-encoded footer image
	BaseImgBase64     string             // Base64-encoded background image
	includeTimestamp  bool               // Whether to include a timestamp
	timestampFormat   string             // Optional timestamp format (default: RFC3339)
	TopRightTimestamp string             // Rendered timestamp value (internal use)
}

// NewPDFRenderer creates a new PDFRenderer.
// The report is any data structure compatible with the given HTML template.
// The images (if used) must be base64-encoded and follow the data URI format.
func NewPDFRenderer(
	report interface{},
	tmpl *template.Template,
	fontSizes renderer.FontSizes,
	useImages bool,
	baseImgBase64, headerImgBase64, footerImgBase64 string,
) interfaces.RendererInterface {
	return &PDFRenderer{
		Report:          report,
		Template:        tmpl,
		FontSizes:       fontSizes,
		UseImages:       useImages,
		HeaderImgBase64: headerImgBase64,
		FooterImgBase64: footerImgBase64,
		BaseImgBase64:   baseImgBase64,
	}
}

// Render generates the PDF output.
// If `filename` is not empty, the output will also be saved to the given file path.
// Returns the generated PDF bytes (in-memory) regardless of file saving.
func (p *PDFRenderer) Render(filename string) ([]byte, error) {
	var buf bytes.Buffer
	if err := p.Template.Execute(&buf, p.Report); err != nil {
		return nil, err
	}

	var r *renderer.Renderer
	var err error
	if p.UseImages {
		r, err = renderer.NewRendererWithBase64Images(p.BaseImgBase64, p.HeaderImgBase64, p.FooterImgBase64, p.FontSizes, true)
	} else {
		r, err = renderer.NewRenderer(p.FontSizes, true)
	}
	if err != nil {
		return nil, err
	}

	if p.includeTimestamp {
		format := p.timestampFormat
		if format == "" {
			format = "2006-01-02 15:04:05"
		}
		r.TopRightTimestamp = time.Now().Format(format)
		log.Println("time:", r.TopRightTimestamp)
	}

	buffer, err := r.RenderHTMLLikeToBuffer(buf.String())
	if err != nil {
		return nil, err
	}

	if filename != "" {
		if err := os.WriteFile(filename, buffer.Bytes(), 0644); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// WithTimestamp enables or disables adding a timestamp in the top-right corner of the PDF pages.
func (p *PDFRenderer) WithTimestamp(enabled bool) interfaces.RendererInterface {
	p.includeTimestamp = enabled
	return p
}

// SetTimestampFormat allows customization of the timestamp format.
// Uses Go's time formatting layout syntax (e.g., "2006-01-02 15:04").
func (p *PDFRenderer) SetTimestampFormat(format string) interfaces.RendererInterface {
	p.timestampFormat = format
	return p
}
