# goreportx

**goreportx** is a lightweight Go library to generate PDF and JSON reports from structured data using HTML-like templates.  
It supports multiple layouts, charts, multi-column tables, and optional timestamp embedding.

---

## Features

- Render structured reports to PDF and JSON
- Embed charts and logos via base64 images
- Support for rich HTML-like templates (`h1`, `h2`, `table`, `img`, `p`, etc.)
- Supports multi-column and chart-based report types
- Optional timestamp injection (top-right for PDF, footer for JSON)
- Modular renderer interface (`RendererInterface`) for extensibility
- Output reports to file or in-memory buffer
- CLI-ready design (planned)

---

## Installation

```bash
go get github.com/ozgen/goreportx
````

---

## Usage

### Render a basic report to PDF and JSON

```go
report := models.Report{
  Header: models.Header{Title: "Smart Report", Subtitle: "Q2", Explanation: "Auto-generated"},
  Footer: models.Footer{Note: "Generated by goreportx"},
  Data: map[string]string{"Project": "AI Dashboard", "Status": "Complete"},
  Charts: []models.Chart{{
    Title: "Usage Overview",
    Tag: renderer.WrapChartAsHTML(chartBase64, renderer.AlignCenter),
    Description: "Chart showing usage trends.",
  }},
}

tmpl, _ := template.ParseFiles("internal/template/defaults/smart_template_new.html")

pdfRenderer := pdf.NewPDFRenderer(
  report,
  tmpl,
  renderer.FontSizes{H1: 24, H2: 20, H3: 16, P: 12, Footer: 10},
  true,
  "", "", "",
).WithTimestamp(true)

pdfRenderer.Render("output.pdf")

jsonRenderer := json.NewJSONRenderer(report).WithTimestamp(true)
jsonRenderer.Render("output.json")
```

---

## Examples

Run example apps from:

```
examples/basic_report/main.go
examples/simple_report/main.go
examples/multiple_column_report/main.go
```

Each example shows how to define templates and models.

---

## Templates

Templates are regular Go `html/template` files using placeholders like:

```html
<h1>{{ .Header.Title }}</h1>
<p>{{ .Footer.Note }}</p>
{{ range .Charts }}{{ .Tag }}{{ end }}
```

---

## Supported Tags in Templates

* `<h1>`, `<h2>`, `<h3>`
* `<p>` (supports `<strong>`, `<em>`)
* `<table>`, `<th>`, `<td>`
* `<img src="data:image/png;base64,...">`
* `<br>` for spacing

---

## Timestamp Support

```go
renderer.WithTimestamp(true).SetTimestampFormat("2006-01-02")
```

* PDF: Top-right of every page
* JSON: Appended to `Footer.Note`

---

## Folder Structure

```
/internal/
  models/         → Report data types
  renderer/       → PDF engine
  renderer/pdf/   → PDFRenderer using template
  renderer/json/  → JSONRenderer
/examples/
  basic_report/   → Charts + tables
  simple_report/  → Logo + single chart
  multiple_column_report/ → Dynamic columns
```

---

## Documentation

* [Project Wiki](https://github.com/ozgen/goreportx/wiki) — Full guide including templates, data models, and examples
* [Templating](https://github.com/ozgen/goreportx/wiki/Templating-in-goreportx)
* [Data Models](https://github.com/ozgen/goreportx/wiki/Data-Models)

---

## TODO

* [x] PDF & JSON renderers
* [x] Templating engine
* [x] Multi-chart support
* [x] Timestamp support
* [ ] CLI interface
* [ ] Unit test
* [ ] Make File support
* [ ] CI/CD pipeline

---

## License

Apache License 2.0 © 2025 [ozgen](https://github.com/ozgen)
