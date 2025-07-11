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

type PDFRenderer struct {
	Report            interface{}
	Template          *template.Template
	FontSizes         renderer.FontSizes
	UseImages         bool
	HeaderImgBase64   string
	FooterImgBase64   string
	BaseImgBase64     string
	includeTimestamp  bool
	timestampFormat   string
	TopRightTimestamp string
}

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

	// Add timestamp to top-right if enabled
	if p.includeTimestamp {
		format := p.timestampFormat
		if format == "" {
			format = "2006-01-02 15:04:05"
		}
		r.TopRightTimestamp = time.Now().Format(format)
		log.Println("time : ", r.TopRightTimestamp)
	}

	buffer, err := r.RenderHTMLLikeToBuffer(buf.String())
	if filename != "" {
		if err := os.WriteFile(filename, buffer.Bytes(), 0644); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func (p *PDFRenderer) WithTimestamp(enabled bool) interfaces.RendererInterface {
	p.includeTimestamp = enabled
	return p
}

func (p *PDFRenderer) SetTimestampFormat(format string) interfaces.RendererInterface {
	p.timestampFormat = format
	return p
}
