package main

import (
	"github.com/ozgen/goreportx/examples/common"
	"github.com/ozgen/goreportx/internal/models"
	"github.com/ozgen/goreportx/internal/renderer"
	"github.com/ozgen/goreportx/internal/renderer/pdf"
	"html/template"
	"log"
	"path/filepath"
)

func main() {
	// Load logo
	logoBase64 := renderer.LoadImageBase64("assets/logo.png")

	// Construct simple report
	report := models.SimpleReport{
		Header: models.SimpleHeader{
			Title:       "Sales Report",
			Subtitle:    "April - June 2025",
			Explanation: "Overview of sales performance for Q2.",
			Logo:        renderer.WrapLogoAsHTML(logoBase64, renderer.AlignCenter),
		},
		Chart: models.SimpleChart{
			Image:       renderer.WrapChartAsHTML(common.GenerateChartBase64(), renderer.AlignLeft),
			Align:       "left",
			Title:       "Quarterly Sales",
			Description: "Sales distribution across regions",
		},
		Footer: models.SimpleFooter{
			Note: "Prepared by GoReportX • All rights reserved.",
		},
	}

	// Load template
	tmplPath := filepath.Join("internal", "template", "defaults", "simple_template.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Template error: %v", err)
	}

	// Render as PDF
	renderer := pdf.NewPDFRenderer(
		report,
		tmpl,
		common.DefaultFontSizes,
		false, // no background/header/footer images
		"",
		"",
		"",
	).WithTimestamp(true)

	_, err = renderer.Render("output-simple.pdf")
	if err != nil {
		log.Fatalf("PDF render error: %v", err)
	}

	log.Println("Simple PDF report generated: output-simple.pdf")
}
