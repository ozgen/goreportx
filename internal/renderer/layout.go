package renderer

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

func (r *Renderer) walk(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1":
			text := GetTextContent(n)
			r.checkPageBreak(30)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H1)
			pageW := 595.28
			textW, _ := r.pdf.MeasureTextWidth(text)
			centerX := (pageW - textW) / 2
			r.pdf.SetX(centerX)
			r.pdf.SetY(r.y)
			r.pdf.Cell(nil, text)
			r.y += 30

		case "h2":
			text := GetTextContent(n)
			r.checkPageBreak(30)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H2)
			r.pdf.SetX(50)
			r.pdf.SetY(r.y)
			r.pdf.Cell(nil, text)
			r.y += 25

		case "h3":
			text := GetTextContent(n)
			r.checkPageBreak(14)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H3)
			r.pdf.SetX(50)
			r.pdf.SetY(r.y)
			r.pdf.Cell(nil, text)
			r.y += 20

		case "p":
			chunks := GetStyledTextChunks(n)
			lineHeight := 16.0
			maxWidth := 495.0

			currentLine := ""
			currentItalic := true
			currentBold := false

			flushLine := func() {
				if currentLine != "" {
					r.checkPageBreak(lineHeight)
					r.pdf.SetX(50)
					r.pdf.SetY(r.y)
					SetFont(r.pdf, currentItalic, currentBold)
					r.pdf.Cell(nil, currentLine)
					r.y += lineHeight
				}
			}

			for _, chunk := range chunks {
				words := strings.Fields(chunk.Text)
				for _, word := range words {
					testLine := strings.TrimSpace(currentLine + " " + word)
					width, _ := r.pdf.MeasureTextWidth(testLine)

					if width > maxWidth || (chunk.Italic != currentItalic || chunk.Bold != currentBold) {
						flushLine()
						currentLine = word
						currentItalic = chunk.Italic
						currentBold = chunk.Bold
					} else {
						if currentLine != "" {
							currentLine += " "
						}
						currentLine += word
					}
				}
			}

			flushLine()
			r.y += 4

		case "table":
			r.renderTable(n)

		case "img":

			var src string
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					src = attr.Val
				}
			}

			// Default alignment
			align := "left"

			// Check parent node's style for alignment
			if n.Parent != nil {
				for _, attr := range n.Parent.Attr {
					if attr.Key == "style" {
						if strings.Contains(attr.Val, "text-align: center") {
							align = "center"
						} else if strings.Contains(attr.Val, "text-align: right") {
							align = "right"
						}
					}
				}
			}

			// Handle base64 or file images
			if strings.HasPrefix(src, "file://") {
				path := strings.TrimPrefix(src, "file://")
				r.drawAlignedImage(path, align)
			} else if strings.HasPrefix(src, "data:image/") {
				if path, err := saveBase64ImageToTempFile(src); err == nil {
					r.drawAlignedImage(path, align)
					os.Remove(path)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.walk(c)
	}
}

func (r *Renderer) drawAlignedImage(path, align string) {
	imgW := 100.0
	imgH := 60.0

	var x float64
	switch align {
	case "center":
		x = (595.28 - imgW) / 2
	case "right":
		x = 595.28 - imgW - 50
	default:
		x = 50
	}

	err := r.pdf.Image(path, x, r.y, &gopdf.Rect{W: imgW, H: imgH})
	if err != nil {
		log.Println("Image render failed:", err)
	} else {
		log.Printf("Image rendered (%s): %s", align, path)
	}
	r.y += imgH + 10
	r.checkPageBreak(30)
}

func (r *Renderer) flushPage() {
	r.drawFooterAtFixedPosition()
	r.pdf.AddPage()
	r.pageNumber++
	r.y = pageMargin

	if r.backgroundImg != "" {
		_ = r.pdf.Image(r.backgroundImg, 0, 0, &gopdf.Rect{W: 595.28, H: 841.89})
	}
	if r.headerImg != "" {
		_ = r.pdf.Image(r.headerImg, 50, 20, &gopdf.Rect{W: 495.0, H: 40.0})
		r.y += 50
	}
	if r.footerImg != "" {
		_ = r.pdf.Image(r.footerImg, 50, 800, &gopdf.Rect{W: 495.0, H: 30.0})
	}
}

func (r *Renderer) drawFooterAtFixedPosition() {
	log.Println("footer:", r.footerText)

	if r.footerText != "" {
		log.Println("footer: %v", r.footerText)

		_ = r.pdf.SetFont("Arial", "", r.FontSize.Footer)
		r.pdf.SetY(820)
		r.pdf.SetX(50)
		r.pdf.Cell(nil, r.footerText)
	}
	if r.showPageNumber {
		_ = r.pdf.SetFont("Arial", "", r.FontSize.Footer)
		r.pdf.SetY(820)
		r.pdf.SetX(500)
		r.pdf.Cell(nil, fmt.Sprintf("Page %d", r.pageNumber))
	}
}

func (r *Renderer) checkPageBreak(nextBlockHeight float64) {
	if r.y+nextBlockHeight > contentLimit {
		log.Println("Page break triggered")
		r.flushPage()
	}
}

func (r *Renderer) renderTable(n *html.Node) {
	_ = r.pdf.SetFont("Arial", "", r.FontSize.P)
	r.walkTableRows(n)
	r.y += 10
}

func (r *Renderer) walkTableRows(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "tr" {
				r.checkPageBreak(30)
				r.renderTableRow(c)
			} else {
				// keep looking inside <thead>, <tbody>, etc.
				r.walkTableRows(c)
			}
		}
	}
}

func (r *Renderer) renderTableRow(tr *html.Node) {
	x := 50.0
	startY := r.y
	rowHeight := 20.0
	colWidth := 250.0

	// Determine if it's a header row
	isHeader := false
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "th" {
			isHeader = true
			break
		}
	}

	if isHeader {
		_ = r.pdf.SetFont("Arial", "B", r.FontSize.P)
	} else {
		_ = r.pdf.SetFont("Arial", "", r.FontSize.P)
	}

	// Draw cells
	for td := tr.FirstChild; td != nil; td = td.NextSibling {
		if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
			text := GetTextContent(td)

			// Draw border
			r.pdf.RectFromUpperLeftWithStyle(x, startY, colWidth, rowHeight, "D")
			r.pdf.SetX(x + 4)
			r.pdf.SetY(startY + 6)
			r.pdf.CellWithOption(&gopdf.Rect{W: colWidth - 8, H: rowHeight}, text, gopdf.CellOption{Align: gopdf.Left})

			x += colWidth
		}
	}

	r.y += rowHeight
}
