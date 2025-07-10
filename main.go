package main

import (
	"bytes"
	"encoding/base64"
	"github.com/ozgen/goreportx/internal/models"
	"github.com/ozgen/goreportx/internal/renderer"
	"github.com/ozgen/goreportx/internal/renderer/json"
	"github.com/ozgen/goreportx/internal/renderer/pdf"
	"github.com/ozgen/goreportx/internal/renderer_old"
	"github.com/wcharczuk/go-chart/v2"
	"html/template"
	"log"
)

func main() {
	align := "center" // or "left", "right" – can be made dynamic

	logoSrc := renderer.LoadImageBase64("assets/logo.png")
	logoTag := template.HTML("")
	if logoSrc != "" {
		logoTag = renderer.WrapLogoAsHTML(logoSrc, align)
	}

	charts := []models.Chart{
		{
			Title:       "Usage Overview",
			Description: "Chart showing daily user activity.",
			Tag:         renderer.WrapChartAsHTML(generateChartBase64(), "center"),
			Order:       0,
		},
		{
			Title:       "Error Trends",
			Description: "Error spikes across regions.",
			Tag:         renderer.WrapChartAsHTML(generateChartBase64(), "right"),
			Order:       1,
		},
		{
			Title:       "Left Chart",
			Description: "Error spikes across regions.",
			Tag:         renderer.WrapChartAsHTML(generateChartBase64(), "left"),
			Order:       2,
		},
	}

	report := models.Report{
		Header: models.Header{
			Title:       "Smart Report",
			Subtitle:    "Auto-generated Example",
			Explanation: "This report is generated using simplified HTML.",
			Logo:        logoTag,
		},
		Footer: models.Footer{
			Note: "Generated with goreportx-lite",
		},
		Data: map[string]string{
			"Customer": "Jane Smith",
			"Email":    "jane@example.com",
			"Project":  "AI Dashboard",
			"Status":   "Complete",
		},
		Charts: charts,
	}

	fontSizes := renderer_old.FontSizes{
		H1:     24.0,
		H2:     18.0,
		H3:     14.0,
		P:      12.0,
		Footer: 10.0,
	}

	// Render template
	tmpl, err := template.ParseFiles("internal/template/defaults/smart_template_new.html")
	if err != nil {
		log.Fatalf("Template error: %v", err)
	}
	log.Println("Template Logo:", report.Header.Logo)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		log.Fatalf("Execute error: %v", err)
	}

	log.Println("buff:", buf.String())

	// Render PDF
	r, err := renderer_old.NewRenderer(fontSizes, true)
	logoBase64 := renderer.LoadImageBase64("assets/header_footer.png")
	ri, err := renderer_old.NewRendererWithBase64Images("", logoBase64, logoBase64, fontSizes, true)

	if err != nil {
		log.Fatalf("PDF init error: %v", err)
	}
	if err := r.RenderHTMLLike(buf.String(), "output.pdf"); err != nil {
		log.Fatalf("Render error: %v", err)
	}
	if err := ri.RenderHTMLLike(buf.String(), "output_image.pdf"); err != nil {
		log.Fatalf("Render error: %v", err)
	}
	log.Println("PDF generated: output.pdf")

	pdfRenderer := pdf.NewPDFRenderer(
		report,
		tmpl,
		renderer.FontSizes{
			H1:     24,
			H2:     20,
			H3:     16,
			P:      12,
			Footer: 10,
		},
		true,
		"",
		logoBase64,
		logoBase64,
	)

	pdfRenderer.Render("output-pdf.pdf")
	log.Println("PDF created successfully: output-pdf.pdf")

	jsonRenderer := json.NewJSONRenderer(report)
	jsonRenderer.Render("output.json")

}

func generateChartBase64() string {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1, 2, 3, 4, 5},
				YValues: []float64{1, 2, 1, 3, 4},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buffer)

	encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return "data:image/png;base64," + encoded
}
