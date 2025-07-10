package renderer_old

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	pageHeight   = 842.0
	pageMargin   = 50.0
	footerHeight = 30.0
	contentLimit = pageHeight - pageMargin - footerHeight
)

type TextChunk struct {
	Text   string
	Italic bool
	Bold   bool
}

type FontPaths struct {
	Regular string
	Bold    string
	Italic  string
}

type FontSizes struct {
	H1     float64
	H2     float64
	H3     float64
	P      float64
	Footer float64
}

type Renderer struct {
	pdf            *gopdf.GoPdf
	y              float64
	pageWidth      float64
	footerText     string
	pageNumber     int
	showPageNumber bool
	backgroundImg  string
	headerImg      string
	footerImg      string
	FontSize       FontSizes
}

func NewRenderer(fontSizes FontSizes, showPageNumber bool) (*Renderer, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	fonts, err := findFontPaths()
	if err != nil {
		return nil, err
	}

	// Register font variants under the same family name: "Arial"
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Regular, gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		return nil, fmt.Errorf("regular font: %w", err)
	}
	log.Println("fonts:", fonts)
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Bold, gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		return nil, fmt.Errorf("bold font: %w", err)
	}
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Italic, gopdf.TtfOption{Style: gopdf.Italic}); err != nil {
		return nil, fmt.Errorf("italic font: %w", err)
	}

	// Set default font
	if err := pdf.SetFont("Arial", "", fontSizes.P); err != nil {
		return nil, err
	}

	return &Renderer{
		pdf:            pdf,
		y:              50,
		pageNumber:     1,
		showPageNumber: showPageNumber,
		FontSize:       fontSizes,
	}, nil
}

func NewRendererWithBase64Images(bgBase64, headerBase64, footerBase64 string, fontSizes FontSizes, showPageNumber bool) (*Renderer, error) {
	bgPath := mustSaveBase64ToTempFile(bgBase64)
	headerPath := mustSaveBase64ToTempFile(headerBase64)
	footerPath := mustSaveBase64ToTempFile(footerBase64)

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	fonts, err := findFontPaths()
	if err != nil {
		return nil, err
	}

	if err := pdf.AddTTFFontWithOption("Arial", fonts.Regular, gopdf.TtfOption{Style: gopdf.Regular}); err != nil {
		return nil, fmt.Errorf("regular font: %w", err)
	}
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Bold, gopdf.TtfOption{Style: gopdf.Bold}); err != nil {
		return nil, fmt.Errorf("bold font: %w", err)
	}
	if err := pdf.AddTTFFontWithOption("Arial", fonts.Italic, gopdf.TtfOption{Style: gopdf.Italic}); err != nil {
		return nil, fmt.Errorf("italic font: %w", err)
	}
	if err := pdf.SetFont("Arial", "", fontSizes.P); err != nil {
		return nil, err
	}

	r := &Renderer{
		pdf:            pdf,
		y:              50,
		pageNumber:     1,
		showPageNumber: showPageNumber,
		backgroundImg:  bgPath,
		headerImg:      headerPath,
		footerImg:      footerPath,
		FontSize:       fontSizes,
	}
	if bgPath != "" {
		_ = r.pdf.Image(bgPath, 0, 0, &gopdf.Rect{W: 595.28, H: 841.89})
	}
	if headerPath != "" {
		_ = r.pdf.Image(headerPath, 50, 20, &gopdf.Rect{W: 495.0, H: 40.0})
		r.y += 50
	}
	if footerPath != "" {
		_ = r.pdf.Image(footerPath, 50, 800, &gopdf.Rect{W: 495.0, H: 30.0})
	}
	return r, nil
}

func mustSaveBase64ToTempFile(dataURI string) string {
	if dataURI == "" {
		return ""
	}
	path, err := saveBase64ImageToTempFile(dataURI)
	if err != nil {
		log.Printf("Failed to save base64 image: %v", err)
		return ""
	}
	return path
}

func (r *Renderer) RenderHTMLLike(content string, output string) error {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract footer content from HTML first
	r.extractFooterText(doc)

	// Walk and render content
	r.walk(doc)

	// Draw footer for the last page
	r.drawFooterAtFixedPosition()

	return r.pdf.WritePdf(output)
}

func (r *Renderer) extractFooterText(n *html.Node) {

	if n.Type == html.ElementNode && n.Data == "div" {
		log.Println(" Footer div found")
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.TrimSpace(attr.Val) == "footer" {
				r.footerText = getTextContent(n)
				log.Println(" Footer found:", r.footerText)
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.extractFooterText(c)
	}
}

func (r *Renderer) walk(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1":
			text := getTextContent(n)
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
			text := getTextContent(n)
			r.checkPageBreak(30)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H2)
			r.pdf.SetX(50)
			r.pdf.SetY(r.y)
			r.pdf.Cell(nil, text)
			r.y += 25

		case "h3":
			text := getTextContent(n)
			r.checkPageBreak(14)
			_ = r.pdf.SetFont("Arial", "", r.FontSize.H3)
			r.pdf.SetX(50)
			r.pdf.SetY(r.y)
			r.pdf.Cell(nil, text)
			r.y += 20

		case "p":
			chunks := getStyledTextChunks(n)
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
					setFont(r.pdf, currentItalic, currentBold)
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
			text := getTextContent(td)

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

func getTextContent(n *html.Node) string {
	var buf bytes.Buffer
	walkText(n, &buf)
	return buf.String()
}

func walkText(n *html.Node, buf io.Writer) {
	if n.Type == html.TextNode {
		buf.Write([]byte(strings.TrimSpace(n.Data)))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkText(c, buf)
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

func saveBase64ImageToTempFile(dataURI string) (string, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URI")
	}

	prefix := parts[0]
	var ext string
	switch {
	case strings.Contains(prefix, "image/png"):
		ext = ".png"
	case strings.Contains(prefix, "image/jpeg"):
		ext = ".jpg"
	case strings.Contains(prefix, "image/svg+xml"):
		ext = ".svg"
	default:
		ext = ".img"
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	tmpfile, err := os.CreateTemp("", "img_*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(decoded); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func getStyledTextChunks(n *html.Node) []TextChunk {
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

func setFont(pdf *gopdf.GoPdf, italic, bold bool) {
	style := ""
	if bold {
		style += "B"
	}
	if italic {
		style += "I"
	}
	_ = pdf.SetFont("Arial", style, 12)
}

func findFontPaths() (FontPaths, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return FontPaths{}, fmt.Errorf("cannot determine caller location")
	}

	// currentFile is the path to this file (e.g., renderer.go)
	baseDir := filepath.Dir(currentFile)
	assetsDir := filepath.Join(baseDir, "..", "..", "assets")

	paths := FontPaths{}

	check := func(filename string) string {
		path := filepath.Join(assetsDir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
		return ""
	}

	paths.Regular = check("LiberationSans-Regular.ttf")
	paths.Bold = check("LiberationSans-Bold.ttf")
	paths.Italic = check("LiberationSans-Italic.ttf")
	return paths, nil
}
