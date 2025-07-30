package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRenderer_CreatesRenderer(t *testing.T) {
	// Skip if fonts not available
	fonts, err := findFontPaths()
	if err != nil || fonts.Regular == "" || fonts.Bold == "" || fonts.Italic == "" {
		t.Skip("Skipping test: fonts not found in assets")
	}

	fontSizes := FontSizes{
		H1: 20, H2: 16, H3: 14, P: 12, Footer: 10,
	}

	renderer, err := NewRenderer(fontSizes, true)

	assert.NoError(t, err)
	assert.NotNil(t, renderer)
	assert.Equal(t, fontSizes, renderer.FontSize)
	assert.True(t, renderer.showPageNumber)
}

func TestNewRendererWithBase64Images_CreatesRendererWithImages(t *testing.T) {
	// Skip if fonts not available
	fonts, err := findFontPaths()
	if err != nil || fonts.Regular == "" || fonts.Bold == "" || fonts.Italic == "" {
		t.Skip("Skipping test: fonts not found in assets")
	}

	// Tiny 1x1 transparent PNG
	tinyPNG := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWNgYGAAAAAEAAGjCh0AAAAASUVORK5CYII="
	base64URI := "data:image/png;base64," + tinyPNG

	fontSizes := FontSizes{
		H1: 24, H2: 18, H3: 14, P: 12, Footer: 10,
	}

	renderer, err := NewRendererWithBase64Images(base64URI, base64URI, base64URI, fontSizes, false)

	assert.NoError(t, err)
	assert.NotNil(t, renderer)
	assert.Equal(t, fontSizes, renderer.FontSize)
	assert.False(t, renderer.showPageNumber)
	assert.NotEmpty(t, renderer.backgroundImg)
	assert.NotEmpty(t, renderer.headerImg)
	assert.NotEmpty(t, renderer.footerImg)
}
