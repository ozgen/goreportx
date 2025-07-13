// File: renderer/html_helpers.go
package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/signintech/gopdf"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Alignment represents possible horizontal alignments for images or HTML elements.
type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

// GetTextContent extracts and returns all visible text from an HTML node and its children.
func GetTextContent(n *html.Node) string {
	var buf bytes.Buffer
	walkText(n, &buf)
	return buf.String()
}

// walkText recursively walks the HTML tree and collects text from TextNode elements.
func walkText(n *html.Node, buf io.Writer) {
	if n.Type == html.TextNode {
		buf.Write([]byte(strings.TrimSpace(n.Data)))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkText(c, buf)
	}
}

// GetStyledTextChunks parses an HTML node tree and returns chunks of text with bold/italic annotations.
func GetStyledTextChunks(n *html.Node) []TextChunk {
	var chunks []TextChunk

	var walk func(node *html.Node, italic, bold bool)
	walk = func(node *html.Node, italic, bold bool) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				chunks = append(chunks, TextChunk{Text: text, Italic: italic, Bold: bold})
			}
		} else if node.Type == html.ElementNode {
			newItalic := italic || node.Data == "em" || node.Data == "i"
			newBold := bold || node.Data == "strong" || node.Data == "b"
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				walk(c, newItalic, newBold)
			}
		}
	}

	walk(n, false, false)
	return chunks
}

// WrapLogoAsHTML returns an HTML <img> tag wrapped in a div with the specified alignment.
func WrapLogoAsHTML(logoBase64 string, align Alignment) template.HTML {
	if logoBase64 == "" {
		return ""
	}
	if align == "" {
		align = AlignLeft
	}
	return template.HTML(fmt.Sprintf(
		`<div style="text-align: %s;"><img src="%s" style="max-height: 60px;" /></div>`,
		align, logoBase64,
	))
}

// WrapChartAsHTML returns a styled HTML <img> tag for a chart, wrapped in a div with alignment.
func WrapChartAsHTML(imgBase64 string, align Alignment) template.HTML {
	if imgBase64 == "" {
		return ""
	}
	if align == "" {
		align = AlignLeft
	}
	return template.HTML(fmt.Sprintf(
		`<div style="text-align: %s;"><img src="%s" style="max-height: 300px;" /></div>`,
		align, imgBase64,
	))
}

// WrapChartAsHTMLWithMeta generates HTML that includes a title, description, and a chart image.
func WrapChartAsHTMLWithMeta(imgBase64, align, title, description string) template.HTML {
	var html string
	if title != "" {
		html += fmt.Sprintf(`<h2 style="text-align: %s;">%s</h2>`, align, title)
	}
	if description != "" {
		html += fmt.Sprintf(`<p style="text-align: %s;">%s</p>`, align, description)
	}
	html += fmt.Sprintf(`<div style="text-align: %s;"><img src="%s" style="max-height: 300px;" /></div>`, align, imgBase64)
	return template.HTML(html)
}

// LoadImageBase64 reads an image from disk and encodes it as a base64 data URI.
func LoadImageBase64(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var mimeType string
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".svg":
		mimeType = "image/svg+xml"
	case ".gif":
		mimeType = "image/gif"
	default:
		log.Fatalf("Unsupported image type: %s", ext)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return "data:" + mimeType + ";base64," + encoded
}

// wrapText splits long text into multiple lines that fit within a given width.
// It returns the lines as a slice of strings.
func wrapText(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	words := strings.Fields(text)
	lines := []string{}
	current := ""

	for _, word := range words {
		test := strings.TrimSpace(current + " " + word)
		width, _ := pdf.MeasureTextWidth(test)
		if width > maxWidth && current != "" {
			lines = append(lines, current)
			current = word
		} else {
			if current != "" {
				current += " "
			}
			current += word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
