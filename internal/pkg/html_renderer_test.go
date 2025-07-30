package pkg

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"testing"
)

func TestRenderHTMLLikeToBuffer_RendersMinimalPDF(t *testing.T) {
	// Skip test if fonts aren't available
	fonts, err := findFontPaths()
	if err != nil || fonts.Regular == "" || fonts.Bold == "" || fonts.Italic == "" {
		t.Skip("Skipping test: fonts not found in assets")
	}

	// Create renderer
	fontSizes := FontSizes{H1: 20, H2: 16, H3: 14, P: 12, Footer: 10}
	renderer, err := NewRenderer(fontSizes, true)
	assert.NoError(t, err)
	assert.NotNil(t, renderer)

	// Sample minimal HTML-like content with a footer
	html := `
		<html>
			<body>
				<h1>Report Title</h1>
				<p>This is a test paragraph.</p>
				<div class="footer">Page 1 of 1</div>
			</body>
		</html>`

	// Act
	buf, err := renderer.RenderHTMLLikeToBuffer(html)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, buf)
	assert.Greater(t, buf.Len(), 100, "PDF buffer should not be empty")

	assert.Equal(t, "Page 1 of 1", renderer.footerText)
}

func TestRenderer_RenderHTMLLikeToBuffer(t *testing.T) {
	type fields struct {
		pdf               *gopdf.GoPdf
		y                 float64
		pageWidth         float64
		footerText        string
		pageNumber        int
		showPageNumber    bool
		backgroundImg     string
		headerImg         string
		footerImg         string
		FontSize          FontSizes
		TopRightTimestamp string
	}
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *bytes.Buffer
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				pdf:               tt.fields.pdf,
				y:                 tt.fields.y,
				pageWidth:         tt.fields.pageWidth,
				footerText:        tt.fields.footerText,
				pageNumber:        tt.fields.pageNumber,
				showPageNumber:    tt.fields.showPageNumber,
				backgroundImg:     tt.fields.backgroundImg,
				headerImg:         tt.fields.headerImg,
				footerImg:         tt.fields.footerImg,
				FontSize:          tt.fields.FontSize,
				TopRightTimestamp: tt.fields.TopRightTimestamp,
			}
			got, err := r.RenderHTMLLikeToBuffer(tt.args.content)
			if !tt.wantErr(t, err, fmt.Sprintf("RenderHTMLLikeToBuffer(%v)", tt.args.content)) {
				return
			}
			assert.Equalf(t, tt.want, got, "RenderHTMLLikeToBuffer(%v)", tt.args.content)
		})
	}
}

func TestRenderer_extractFooterText(t *testing.T) {
	type fields struct {
		pdf               *gopdf.GoPdf
		y                 float64
		pageWidth         float64
		footerText        string
		pageNumber        int
		showPageNumber    bool
		backgroundImg     string
		headerImg         string
		footerImg         string
		FontSize          FontSizes
		TopRightTimestamp string
	}
	type args struct {
		n *html.Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				pdf:               tt.fields.pdf,
				y:                 tt.fields.y,
				pageWidth:         tt.fields.pageWidth,
				footerText:        tt.fields.footerText,
				pageNumber:        tt.fields.pageNumber,
				showPageNumber:    tt.fields.showPageNumber,
				backgroundImg:     tt.fields.backgroundImg,
				headerImg:         tt.fields.headerImg,
				footerImg:         tt.fields.footerImg,
				FontSize:          tt.fields.FontSize,
				TopRightTimestamp: tt.fields.TopRightTimestamp,
			}
			r.extractFooterText(tt.args.n)
		})
	}
}
