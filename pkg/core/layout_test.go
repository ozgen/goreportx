package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func defaultFontSizes() FontSizes {
	return FontSizes{H1: 24, H2: 18, H3: 14, P: 12, Footer: 10}
}

func fontsMissing() bool {
	fonts, err := findFontPaths()
	return err != nil || fonts.Regular == "" || fonts.Bold == "" || fonts.Italic == ""
}

func TestRenderHTMLLikeToBuffer_H1Header(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found, skipping")
	}

	r, _ := NewRenderer(defaultFontSizes(), true)
	html := `<h1>Big Title</h1>`

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
}

func TestRenderHTMLLikeToBuffer_ParagraphStyled(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found")
	}

	r, _ := NewRenderer(defaultFontSizes(), true)
	html := `<p>This is <strong>bold</strong> and <em>italic</em>.</p>`

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
}

func TestRenderHTMLLikeToBuffer_MultipleHeaders(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found")
	}
	r, _ := NewRenderer(defaultFontSizes(), false)
	html := `<h1>Main</h1><h2>Sub</h2><h3>Minor</h3>`

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
}

func TestRenderHTMLLikeToBuffer_InlineImage(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found")
	}
	r, _ := NewRenderer(defaultFontSizes(), false)

	img := `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWNgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII=`
	html := fmt.Sprintf(`<div style="text-align: center;"><img src="%s" /></div>`, img)

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
}

func TestRenderHTMLLikeToBuffer_SimpleTable(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found")
	}
	r, _ := NewRenderer(defaultFontSizes(), true)

	html := `<table>
	<tr><th>Header 1</th><th>Header 2</th></tr>
	<tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
	<tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
	</table>`

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Greater(t, buf.Len(), 100)
}

func TestRenderHTMLLikeToBuffer_FooterAndTimestamp(t *testing.T) {
	if fontsMissing() {
		t.Skip("Fonts not found")
	}
	r, _ := NewRenderer(defaultFontSizes(), true)
	r.TopRightTimestamp = "July 2025"

	html := `<div class="footer">Confidential Report</div><p>Content.</p>`

	buf, err := r.RenderHTMLLikeToBuffer(html)
	assert.NoError(t, err)
	assert.Equal(t, "Confidential Report", r.footerText)
	assert.Greater(t, buf.Len(), 100)
}
