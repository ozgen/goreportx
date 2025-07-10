package json

import (
	"encoding/json"
	"os"

	"github.com/ozgen/goreportx/internal/interfaces"
	"github.com/ozgen/goreportx/internal/models"
)

type JSONRenderer struct {
	Report models.Report
}

func NewJSONRenderer(report models.Report) interfaces.RendererInterface {
	return &JSONRenderer{Report: report}
}

func (r *JSONRenderer) Render(filename string) ([]byte, error) {
	data, err := json.MarshalIndent(r.Report, "", "  ")
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
