package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/ozgen/goreportx/internal/core"
	jsonReort "github.com/ozgen/goreportx/internal/renderer/json"
	"github.com/ozgen/goreportx/internal/renderer/pdf"
)

func main() {
	// Required
	inputPath := flag.String("input", "", "Path to JSON input file (required)")
	templatePath := flag.String("template", "", "Path to HTML template file (required)")
	format := flag.String("format", "", "Output format: pdf or json (required)")

	// Optional
	outputPath := flag.String("output", "", "Output file path")
	headerImg := flag.String("headerImage", "", "Header image path (optional)")
	footerImg := flag.String("footerImage", "", "Footer image path (optional)")
	baseImg := flag.String("baseImage", "", "Background image path (optional)")
	showPageNumber := flag.Bool("showPageNumber", true, "Show page numbers (PDF only)")
	withTimestamp := flag.Bool("with-timestamp", false, "Include timestamp")
	timeFormat := flag.String("time-format", "2006-01-02 15:04:05", "Timestamp format")

	flag.Parse()

	// Validate required flags
	if *inputPath == "" || *templatePath == "" || *format == "" {
		log.Fatal("Missing required flag(s): --input, --template, and --format are mandatory")
	}

	// Default output
	if *outputPath == "" {
		if *format == "pdf" {
			*outputPath = "output.pdf"
		} else if *format == "json" {
			*outputPath = "output.json"
		}
	}

	// Load report JSON
	jsonBytes, err := os.ReadFile(*inputPath)
	if err != nil {
		log.Fatalf("Failed to read input JSON: %v", err)
	}
	var report map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &report); err != nil {
		log.Fatalf("Failed to parse input JSON: %v", err)
	}

	// Load template
	tmpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Choose format
	switch strings.ToLower(*format) {
	case "pdf":
		factory := core.NewRendererFactory().WithPageNumbers(*showPageNumber)

		// Optional images
		if *headerImg != "" {
			factory.WithHeaderImage(core.LoadImageBase64(*headerImg))
		}
		if *footerImg != "" {
			factory.WithFooterImage(core.LoadImageBase64(*footerImg))
		}
		if *baseImg != "" {
			factory.WithBaseImage(core.LoadImageBase64(*baseImg))
		}

		// Create PDF renderer
		pdfRenderer := pdf.NewPDFRenderer(report, tmpl, factory).
			WithTimestamp(*withTimestamp).
			SetTimestampFormat(*timeFormat)

		_, err := pdfRenderer.Render(*outputPath)
		if err != nil {
			log.Fatalf("Failed to render PDF: %v", err)
		}
		log.Println("PDF generated at:", *outputPath)

	case "json":
		jsonRenderer := jsonReort.NewJSONRenderer(report).
			WithTimestamp(*withTimestamp).
			SetTimestampFormat(*timeFormat)

		_, err := jsonRenderer.Render(*outputPath)
		if err != nil {
			log.Fatalf("Failed to render JSON: %v", err)
		}
		log.Println("JSON saved at:", *outputPath)

	default:
		log.Fatalf("Invalid format: %s. Use 'pdf' or 'json'.", *format)
	}
}
