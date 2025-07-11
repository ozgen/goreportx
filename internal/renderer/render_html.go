package renderer

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func (r *Renderer) RenderHTMLLike(content string, output string) error {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	r.extractFooterText(doc)
	r.walk(doc)
	r.drawFooterAtFixedPosition()
	r.drawTimestamp()

	return r.pdf.WritePdf(output)
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
