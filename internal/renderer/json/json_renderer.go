// File: renderer/json/json_renderer.go
package json

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ozgen/goreportx/internal/interfaces"
)

// JSONRenderer is responsible for rendering any report object into
// a structured JSON output. Optionally, a timestamp can be included.
type JSONRenderer struct {
	Report           interface{} // The report data (can be any structure matching the template)
	includeTimestamp bool        // If true, appends timestamp to the Footer.Note
	timestampFormat  string      // Optional custom format (Go layout style)
}

// NewJSONRenderer creates a new JSONRenderer instance.
func NewJSONRenderer(report interface{}) interfaces.RendererInterface {
	return &JSONRenderer{Report: report}
}

// Render marshals the report into indented JSON.
// If a filename is given, it saves the file; otherwise, it returns the JSON as bytes.
// If includeTimestamp is enabled, a timestamp is appended to Footer.Note.
func (r *JSONRenderer) Render(filename string) ([]byte, error) {
	var report interface{} = r.Report

	if r.includeTimestamp {
		// Re-encode into a mutable map to inject timestamp
		jsonBytes, _ := json.Marshal(r.Report)
		var generic map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &generic); err == nil {
			if footer, ok := generic["Footer"].(map[string]interface{}); ok {
				note, _ := footer["Note"].(string)

				format := r.timestampFormat
				if format == "" {
					format = "2006-01-02 15:04:05" // Default Go format
				}

				footer["Note"] = note + " | Generated at: " + time.Now().Format(format)
				generic["Footer"] = footer
			}
			r.Report = generic
			report = generic
		}
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, err
	}

	if filename != "" {
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// WithTimestamp enables or disables appending a timestamp to Footer.Note.
func (r *JSONRenderer) WithTimestamp(enable bool) interfaces.RendererInterface {
	r.includeTimestamp = enable
	return r
}

// SetTimestampFormat allows setting a custom Go-style timestamp format.
// Example: "02 Jan 2006 15:04" or time.RFC3339.
func (r *JSONRenderer) SetTimestampFormat(format string) interfaces.RendererInterface {
	r.timestampFormat = format
	return r
}
