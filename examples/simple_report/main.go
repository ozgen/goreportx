package main

import (
	"github.com/ozgen/goreportx/examples/common"
	"github.com/ozgen/goreportx/internal/core"
	"github.com/ozgen/goreportx/internal/models"
	"github.com/ozgen/goreportx/internal/renderer/pdf"
	"html/template"
	"log"
	"path/filepath"
)

func main() {
	// Load logo
	logoBase64 := core.LoadImageBase64("assets/logo.png")
	logoHTML := core.WrapLogoAsHTML(logoBase64, core.AlignCenter)

	// Construct report model
	report := models.SimpleReport{
		Header: models.SimpleHeader{
			Title:       "Sales Report",
			Subtitle:    "April - June 2025",
			Explanation: "Overview of sales performance for Q2.",
			Logo:        logoHTML,
		},
		Chart: models.SimpleChart{
			Image:       core.WrapChartAsHTML(common.GenerateChartBase64(), core.AlignLeft),
			Align:       "left",
			Title:       "Quarterly Sales",
			Description: "Sales distribution across regions",
		},
		Footer: models.SimpleFooter{
			Note: "Prepared by GoReportX • All rights reserved.",
		},
	}

	// Load HTML template
	tmplPath := filepath.Join("internal", "template", "defaults", "simple_template.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Build renderer factory

	factory := core.NewRendererFactory().
		WithFontSizes(common.DefaultFontSizes).
		WithPageNumbers(false)

	// Create PDF renderer with factory and timestamp
	renderer := pdf.NewPDFRenderer(report, tmpl, factory).
		WithTimestamp(true).
		SetTimestampFormat("2006-01-02 15:04")

	_, err = renderer.Render("examples/outputs/output-simple.pdf")
	if err != nil {
		log.Fatalf("Failed to render PDF: %v", err)
	}

	log.Println("Simple PDF report generated: examples/outputs/output-simple.pdf")
}
