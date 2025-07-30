package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRendererFactory_Defaults(t *testing.T) {
	factory := NewRendererFactory()

	assert.Equal(t, FontSizes{
		H1:     24,
		H2:     18,
		H3:     14,
		P:      12,
		Footer: 10,
	}, factory.FontSizes)

	assert.True(t, factory.ShowPageNumber)
	assert.Empty(t, factory.Base64Background)
	assert.Empty(t, factory.Base64Header)
	assert.Empty(t, factory.Base64Footer)
}

func TestRendererFactory_WithFontSizes(t *testing.T) {
	customSizes := FontSizes{
		H1:     30,
		H2:     22,
		H3:     16,
		P:      13,
		Footer: 11,
	}
	factory := NewRendererFactory().WithFontSizes(customSizes)

	assert.Equal(t, customSizes, factory.FontSizes)
}

func TestRendererFactory_WithPageNumbers(t *testing.T) {
	factory := NewRendererFactory().WithPageNumbers(false)
	assert.False(t, factory.ShowPageNumber)
}

func TestRendererFactory_WithImages(t *testing.T) {
	bg := "base64-bg"
	hdr := "base64-hdr"
	ftr := "base64-ftr"

	factory := NewRendererFactory().
		WithBaseImage(bg).
		WithHeaderImage(hdr).
		WithFooterImage(ftr)

	assert.Equal(t, bg, factory.Base64Background)
	assert.Equal(t, hdr, factory.Base64Header)
	assert.Equal(t, ftr, factory.Base64Footer)
}

func TestRendererFactory_Build_DefaultRenderer(t *testing.T) {
	factory := NewRendererFactory()

	renderer, err := factory.Build()
	assert.NoError(t, err)
	assert.NotNil(t, renderer)

	assert.Equal(t, factory.FontSizes, renderer.FontSize)
	assert.Equal(t, factory.ShowPageNumber, renderer.showPageNumber)
}

func TestRendererFactory_Build_WithBase64Images(t *testing.T) {
	factory := NewRendererFactory().
		WithBaseImage("bg").
		WithHeaderImage("hdr").
		WithFooterImage("ftr")

	renderer, err := factory.Build()
	assert.NoError(t, err)
	assert.NotNil(t, renderer)

	assert.Equal(t, factory.FontSizes, renderer.FontSize)
	assert.Equal(t, factory.ShowPageNumber, renderer.showPageNumber)
}
