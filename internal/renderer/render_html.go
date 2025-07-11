package renderer

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

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
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("failed to remove temp file: %s\n", name)
		}
	}(tmpFile.Name())
	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {
			fmt.Printf("failed to close temp file: %s\n", tmpFile.Name())
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
