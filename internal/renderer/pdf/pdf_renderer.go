// File: internal/renderer/pdf/pdf_renderer.go
package pdf

import (
	"bytes"
	"github.com/ozgen/goreportx/internal/core"
	"github.com/ozgen/goreportx/internal/interfaces"
	"html/template"
	"os"
	"time"
)

// PDFRenderer renders a report using an HTML template into a PDF.
// It uses RendererFactory to configure the internal PDF engine.
type PDFRenderer struct {
	Report            interface{}
	Template          *template.Template
	Factory           *core.RendererFactory
	includeTimestamp  bool
	timestampFormat   string
	TopRightTimestamp string
}

// NewPDFRenderer constructs a PDFRenderer using a given RendererFactory.
// This removes tight coupling with image setup and font handling.
func NewPDFRenderer(
	report interface{},
	tmpl *template.Template,
	factory *core.RendererFactory,
) interfaces.RendererInterface {
	return &PDFRenderer{
		Report:   report,
		Template: tmpl,
		Factory:  factory,
	}
}

// Render generates the PDF and optionally writes to disk.
func (p *PDFRenderer) Render(filename string) ([]byte, error) {
	var buf bytes.Buffer
	if err := p.Template.Execute(&buf, p.Report); err != nil {
		return nil, err
	}

	renderer, err := p.Factory.Build()
	if err != nil {
		return nil, err
	}

	if p.includeTimestamp {
		format := p.timestampFormat
		if format == "" {
			format = "2006-01-02 15:04:05"
		}
		renderer.TopRightTimestamp = time.Now().Format(format)
	}

	pdfBuffer, err := renderer.RenderHTMLLikeToBuffer(buf.String())
	if err != nil {
		return nil, err
	}

	if filename != "" {
		if err := os.WriteFile(filename, pdfBuffer.Bytes(), 0644); err != nil {
			return nil, err
		}
	}

	return pdfBuffer.Bytes(), nil
}

// WithTimestamp enables/disables timestamp display in top-right.
func (p *PDFRenderer) WithTimestamp(enable bool) interfaces.RendererInterface {
	p.includeTimestamp = enable
	return p
}

// SetTimestampFormat specifies custom layout format for timestamp.
func (p *PDFRenderer) SetTimestampFormat(format string) interfaces.RendererInterface {
	p.timestampFormat = format
	return p
}
