package json

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ozgen/goreportx/internal/interfaces"
)

type JSONRenderer struct {
	Report           interface{}
	includeTimestamp bool
	timestampFormat  string
}

func NewJSONRenderer(report interface{}) interfaces.RendererInterface {
	return &JSONRenderer{Report: report}
}

func (r *JSONRenderer) Render(filename string) ([]byte, error) {
	var report interface{} = r.Report

	if r.includeTimestamp {
		jsonBytes, _ := json.Marshal(r.Report)
		var generic map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &generic); err == nil {
			if footer, ok := generic["Footer"].(map[string]interface{}); ok {
				note, _ := footer["Note"].(string)

				format := r.timestampFormat
				if format == "" {
					format = "2006-01-02 15:04:05" // default
				}

				footer["Note"] = note + " | Generated at: " + time.Now().Format(format)
				generic["Footer"] = footer
			}
			r.Report = generic
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

func (r *JSONRenderer) WithTimestamp(enable bool) interfaces.RendererInterface {
	r.includeTimestamp = enable
	return r
}

func (r *JSONRenderer) SetTimestampFormat(format string) interfaces.RendererInterface {
	r.timestampFormat = format
	return r
}
