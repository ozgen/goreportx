package renderer

import (
	"fmt"
	"github.com/signintech/gopdf"
)

func NewRenderer(fontSizes FontSizes, showPageNumber bool) (*Renderer, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	fonts, err := findFontPaths()
	if err != nil {
		return nil, err
	}

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
