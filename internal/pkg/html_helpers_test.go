package pkg

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/signintech/gopdf"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestGetTextContent_ExtractsText(t *testing.T) {
	htmlStr := `<div><p>Hello <strong>world</strong>!</p></div>`
	node, _ := html.Parse(strings.NewReader(htmlStr))

	text := GetTextContent(node)
	assert.Contains(t, text, "Hello")
	assert.Contains(t, text, "world")
}

func TestGetStyledTextChunks_BoldItalic(t *testing.T) {
	htmlStr := `<p>This is <strong>bold</strong> and <em>italic</em> text.</p>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))

	// Find the <p> node
	var pNode *html.Node
	var findP func(*html.Node)
	findP = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			pNode = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findP(c)
		}
	}
	findP(doc)
	assert.NotNil(t, pNode, "expected to find <p> node")

	chunks := GetStyledTextChunks(pNode)
	assert.GreaterOrEqual(t, len(chunks), 4)

	// Optionally check exact text/style
	assert.Equal(t, "This is", chunks[0].Text)
	assert.False(t, chunks[0].Bold)
	assert.False(t, chunks[0].Italic)

	assert.Equal(t, "bold", chunks[1].Text)
	assert.True(t, chunks[1].Bold)

	assert.Equal(t, "and", chunks[2].Text)
	assert.False(t, chunks[2].Bold)

	assert.Equal(t, "italic", chunks[3].Text)
	assert.True(t, chunks[3].Italic)
}

func TestWrapLogoAsHTML_GeneratesValidHTML(t *testing.T) {
	html := WrapLogoAsHTML("data:image/png;base64,fake", AlignCenter)
	assert.Contains(t, string(html), "text-align: center")
	assert.Contains(t, string(html), "<img src=")
}

func TestWrapChartAsHTML_GeneratesValidHTML(t *testing.T) {
	html := WrapChartAsHTML("data:image/png;base64,xyz", AlignRight)
	assert.Contains(t, string(html), "text-align: right")
	assert.Contains(t, string(html), `<img src="data:image/png;base64,xyz"`)
}

func TestWrapLogoAsHTML_DefaultsToAlignLeft(t *testing.T) {
	base64 := "data:image/png;base64,someFakeImage"
	html := WrapLogoAsHTML(base64, "")

	assert.Contains(t, string(html), `text-align: left`, "Expected default alignment to be 'left'")
	assert.Contains(t, string(html), `<img src="`+base64+`"`)
	assert.Contains(t, string(html), `max-height: 60px;`)
}

func TestWrapChartAsHTMLWithMeta_CombinesElements(t *testing.T) {
	html := WrapChartAsHTMLWithMeta("data:image/png;base64,123", "center", "Sales", "Monthly report")
	out := string(html)

	assert.Contains(t, out, "<h2 style=\"text-align: center;\">Sales</h2>")
	assert.Contains(t, out, "<p style=\"text-align: center;\">Monthly report</p>")
	assert.Contains(t, out, "data:image/png;base64,123")
}

func TestWrapText_BreaksLongText(t *testing.T) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	_ = pdf.AddTTFFont("Arial", "../../assets/LiberationSans-Regular.ttf") // Skip error check for now
	_ = pdf.SetFont("Arial", "", 12)

	text := "This is a very long line of text that should break into multiple lines based on width."
	lines := wrapText(pdf, text, 100) // Narrow width

	assert.Greater(t, len(lines), 1)
	for _, line := range lines {
		assert.NotEmpty(t, line)
	}
}

func TestLoadImageBase64_PNG(t *testing.T) {
	// Arrange: Create a minimal PNG file
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")

	// This is a 1x1 transparent PNG
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x60, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x01, 0xE5, 0x27, 0xD4, 0xA2, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	err := os.WriteFile(imgPath, pngBytes, 0644)
	assert.NoError(t, err)

	// Act
	base64Str := LoadImageBase64(imgPath)

	// Assert
	assert.True(t, strings.HasPrefix(base64Str, "data:image/png;base64,"), "should contain correct MIME prefix")
	assert.Greater(t, len(base64Str), 30, "base64 output should not be empty")
}

func TestLoadImageBase64_MissingFilePanics(t *testing.T) {
	if os.Getenv("TEST_PANIC_MISSING_FILE") == "1" {
		_ = LoadImageBase64("non_existent_file.png")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoadImageBase64_MissingFilePanics")
	cmd.Env = append(os.Environ(), "TEST_PANIC_MISSING_FILE=1")
	output, err := cmd.CombinedOutput()

	assert.Error(t, err)
	assert.Contains(t, string(output), "Failed to read image")
}

func TestLoadImageBase64_UnsupportedExtensionPanics(t *testing.T) {
	tmpDir := t.TempDir()
	badPath := filepath.Join(tmpDir, "bad.bmp")
	err := os.WriteFile(badPath, []byte("dummy"), 0644)
	assert.NoError(t, err)

	if os.Getenv("TEST_PANIC_UNSUPPORTED_EXT") == "1" {
		_ = LoadImageBase64(badPath)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoadImageBase64_UnsupportedExtensionPanics")
	cmd.Env = append(os.Environ(), "TEST_PANIC_UNSUPPORTED_EXT=1")
	output, err := cmd.CombinedOutput()

	assert.Error(t, err)
	assert.Contains(t, string(output), "Unsupported image type")
}
