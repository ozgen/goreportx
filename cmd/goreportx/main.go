package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/ozgen/goreportx/internal/core"
	jsonReort "github.com/ozgen/goreportx/internal/renderer/json"
	"github.com/ozgen/goreportx/internal/renderer/pdf"
)

func Run(args []string) error {
	fs := flag.NewFlagSet("goreportx", flag.ContinueOnError)

	inputPath := fs.String("input", "", "Path to JSON input file (required)")
	templatePath := fs.String("template", "", "Path to HTML template file (required)")
	format := fs.String("format", "", "Output format: pdf or json (required)")
	outputPath := fs.String("output", "", "Output file path")
	headerImg := fs.String("headerImage", "", "Header image path (optional)")
	footerImg := fs.String("footerImage", "", "Footer image path (optional)")
	baseImg := fs.String("baseImage", "", "Background image path (optional)")
	showPageNumber := fs.Bool("showPageNumber", true, "Show page numbers (PDF only)")
	withTimestamp := fs.Bool("with-timestamp", false, "Include timestamp")
	timeFormat := fs.String("time-format", "2006-01-02 15:04:05", "Timestamp format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *inputPath == "" || *templatePath == "" || *format == "" {
		return errors.New("missing required flag(s): --input, --template, and --format are mandatory")
	}

	if *outputPath == "" {
		if *format == "pdf" {
			*outputPath = "output.pdf"
		} else if *format == "json" {
			*outputPath = "output.json"
		}
	}

	jsonBytes, err := os.ReadFile(*inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input JSON: %w", err)
	}
	var report map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &report); err != nil {
		return fmt.Errorf("failed to parse input JSON: %w", err)
	}

	tmpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	switch strings.ToLower(*format) {
	case "pdf":
		factory := core.NewRendererFactory().WithPageNumbers(*showPageNumber)
		if *headerImg != "" {
			factory.WithHeaderImage(core.LoadImageBase64(*headerImg))
		}
		if *footerImg != "" {
			factory.WithFooterImage(core.LoadImageBase64(*footerImg))
		}
		if *baseImg != "" {
			factory.WithBaseImage(core.LoadImageBase64(*baseImg))
		}

		_, err := pdf.NewPDFRenderer(report, tmpl, factory).
			WithTimestamp(*withTimestamp).
			SetTimestampFormat(*timeFormat).
			Render(*outputPath)
		if err != nil {
			return fmt.Errorf("failed to render PDF: %w", err)
		}

	case "json":
		_, err := jsonReort.NewJSONRenderer(report).
			WithTimestamp(*withTimestamp).
			SetTimestampFormat(*timeFormat).
			Render(*outputPath)
		if err != nil {
			return fmt.Errorf("failed to render JSON: %w", err)
		}

	default:
		return fmt.Errorf("invalid format: %s. Use 'pdf' or 'json'", *format)
	}

	return nil
}

// main() just calls Run with os.Args[1:]
func main() {
	if err := Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
