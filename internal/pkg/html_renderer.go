// File: renderer/html_renderer.go
package pkg

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

// RenderHTMLLikeToBuffer renders simplified HTML-like content to a PDF,
// writes it to a temporary file, reads the file into memory, and returns it as a byte buffer.
// This is useful for cases where you need the PDF in memory (e.g., HTTP response, tests)
// instead of saving it to disk.
func (r *Renderer) RenderHTMLLikeToBuffer(content string) (*bytes.Buffer, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	r.extractFooterText(doc)
	r.walk(doc)
	r.drawFooterAtFixedPosition()
	r.drawTimestamp()

	// Write to a temporary file
	tmpFile, err := os.CreateTemp("", "report_*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			fmt.Printf("failed to remove temp file: %s\n", name)
		}
	}(tmpFile.Name())
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close temp file: %s\n", f.Name())
		}
	}(tmpFile)

	if err := r.pdf.WritePdf(tmpFile.Name()); err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	// Read the contents back into a buffer
	pdfBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read temp PDF: %w", err)
	}

	return bytes.NewBuffer(pdfBytes), nil
}

// extractFooterText looks for a <div class="footer"> element in the HTML document
// and extracts its text content to be used as the page footer.
func (r *Renderer) extractFooterText(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.TrimSpace(attr.Val) == "footer" {
				r.footerText = GetTextContent(n)
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.extractFooterText(c)
	}
}
