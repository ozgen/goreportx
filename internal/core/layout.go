// File: renderer/layout.go
package core

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"

	"github.com/signintech/gopdf"
)

// walk recursively traverses an HTML node tree and renders
// supported elements such as h1–h3, p, table, img, and br to the PDF.
func (r *Renderer) walk(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1":
			text := GetTextContent(n)
			r.checkPageBreak(30)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H1)
			textW, _ := r.pdf.MeasureTextWidth(text)
			r.pdf.SetX((595.28 - textW) / 2)
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
			r.renderParagraph(n)

		case "table":
			r.renderTable(n)

		case "br":
			r.y += 10 // handle line breaks with vertical space

		case "img":
			r.renderImage(n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.walk(c)
	}
}

// renderParagraph processes a <p> element and wraps styled text into lines.
func (r *Renderer) renderParagraph(n *html.Node) {
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

			if width > maxWidth || chunk.Italic != currentItalic || chunk.Bold != currentBold {
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
}

// renderImage processes an <img> element, extracting alignment and base64/file path.
func (r *Renderer) renderImage(n *html.Node) {
	var src string
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			src = attr.Val
			break
		}
	}
	align := "left"
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

	if strings.HasPrefix(src, "file://") {
		path := strings.TrimPrefix(src, "file://")
		r.drawAlignedImage(path, align)
	} else if strings.HasPrefix(src, "data:image/") {
		if path, err := saveBase64ImageToTempFile(src); err == nil {
			r.drawAlignedImage(path, align)
			err := os.Remove(path)
			if err != nil {
				return
			}
		}
	}
}

// drawAlignedImage renders an image with specified horizontal alignment.
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

	if err := r.pdf.Image(path, x, r.y, &gopdf.Rect{W: imgW, H: imgH}); err != nil {
		log.Println("Image render failed:", err)
	} else {
		log.Printf("Image rendered (%s): %s", align, path)
	}
	r.y += imgH + 10
	r.checkPageBreak(30)
}

// flushPage finishes the current page, adds footer/timestamp, and starts a new page.
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

	r.drawTimestamp()
}

// drawFooterAtFixedPosition draws static footer and page number at the bottom of each page.
func (r *Renderer) drawFooterAtFixedPosition() {
	if r.footerText != "" {
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

// checkPageBreak checks if the current y-position plus upcoming block height
// will overflow the page, and triggers a flushPage if so.
func (r *Renderer) checkPageBreak(nextBlockHeight float64) {
	if r.y+nextBlockHeight > contentLimit {
		log.Println("Page break triggered")
		r.flushPage()
	}
}

// renderTable walks through the table node and renders its rows.
func (r *Renderer) renderTable(n *html.Node) {
	_ = r.pdf.SetFont("Arial", "", r.FontSize.P)
	r.walkTableRows(n)
	r.y += 10
}

// walkTableRows traverses table rows and renders each.
func (r *Renderer) walkTableRows(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "tr" {
				r.checkPageBreak(30)
				r.renderTableRow(c)
			} else {
				r.walkTableRows(c)
			}
		}
	}
}

// renderTableRow renders a single table row with dynamic height and column width.
func (r *Renderer) renderTableRow(tr *html.Node) {
	x := 50.0
	lineHeight := 14.0
	numCols := r.countColumns(tr)
	if numCols == 0 {
		return
	}
	colWidth := (595.28 - 100) / float64(numCols)

	maxLines := 1
	cellTexts := []string{}
	for td := tr.FirstChild; td != nil; td = td.NextSibling {
		if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
			text := GetTextContent(td)
			cellTexts = append(cellTexts, text)
			lines := wrapText(r.pdf, text, colWidth-8)
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}
	}
	rowHeight := float64(maxLines) * lineHeight
	r.checkPageBreak(rowHeight)

	isHeader := tr.FirstChild != nil && tr.FirstChild.Data == "th"
	if isHeader {
		_ = r.pdf.SetFont("Arial", "B", r.FontSize.P)
	} else {
		_ = r.pdf.SetFont("Arial", "", r.FontSize.P)
	}

	x = 50.0
	startY := r.y
	for _, text := range cellTexts {
		r.pdf.RectFromUpperLeftWithStyle(x, startY, colWidth, rowHeight, "D")
		lines := wrapText(r.pdf, text, colWidth-8)
		for j, line := range lines {
			r.pdf.SetX(x + 4)
			r.pdf.SetY(startY + float64(j)*lineHeight + 2)
			r.pdf.Cell(nil, line)
		}
		x += colWidth
	}
	r.y += rowHeight
}

// countColumns counts how many <td> or <th> elements are in a given <tr>.
func (r *Renderer) countColumns(tr *html.Node) int {
	count := 0
	for td := tr.FirstChild; td != nil; td = td.NextSibling {
		if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
			count++
		}
	}
	return count
}

// drawTimestamp renders a timestamp string in the top-right of the page if set.
func (r *Renderer) drawTimestamp() {
	if r.TopRightTimestamp == "" {
		return
	}
	_ = r.pdf.SetFont("Arial", "", r.FontSize.Footer)
	r.pdf.SetX(490)
	r.pdf.SetY(30)
	r.pdf.Cell(nil, r.TopRightTimestamp)
}
