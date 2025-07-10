package pdf

import (
	"bytes"
	"github.com/ozgen/goreportx/internal/interfaces"
	"github.com/ozgen/goreportx/internal/models"
	"github.com/ozgen/goreportx/internal/renderer"
	"html/template"
)

type PDFRenderer struct {
	Report          models.Report
	Template        *template.Template
	FontSizes       renderer.FontSizes
	UseImages       bool
	HeaderImgBase64 string
	FooterImgBase64 string
	BaseImgBase64   string
}

func NewPDFRenderer(
	report models.Report,
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

	// 3. Create in-memory PDF
	var outBuf bytes.Buffer
	if err := r.RenderHTMLLike(buf.String(), filename); err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil
}
