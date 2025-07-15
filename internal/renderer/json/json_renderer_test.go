package json

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ozgen/goreportx/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func sampleReport() map[string]interface{} {
	return map[string]interface{}{
		"Title": "Sample Report",
		"Footer": map[string]interface{}{
			"Note": "Original Note",
		},
	}
}

func TestJSONRenderer_Render_Basic(t *testing.T) {
	renderer := NewJSONRenderer(sampleReport())
	result, err := renderer.Render("")

	assert.NoError(t, err)
	assert.Contains(t, string(result), "Sample Report")
	assert.Contains(t, string(result), "Original Note")
}

func TestJSONRenderer_Render_WithTimestamp(t *testing.T) {
	renderer := NewJSONRenderer(sampleReport()).
		WithTimestamp(true)

	result, err := renderer.Render("")
	assert.NoError(t, err)

	var parsed map[string]interface{}
	_ = json.Unmarshal(result, &parsed)

	footer := parsed["Footer"].(map[string]interface{})
	assert.Contains(t, footer["Note"].(string), "Generated at:")
}

func TestJSONRenderer_Render_WithCustomTimestampFormat(t *testing.T) {
	renderer := NewJSONRenderer(sampleReport()).
		WithTimestamp(true).
		SetTimestampFormat("02-Jan-2006 15:04")

	result, err := renderer.Render("")
	assert.NoError(t, err)

	assert.Contains(t, string(result), "Generated at: ")
	assert.Contains(t, string(result), time.Now().Format("02-Jan-2006"))
}

func TestJSONRenderer_Render_WriteToFile(t *testing.T) {
	tmpFile := "test_report.json"
	defer os.Remove(tmpFile)

	renderer := NewJSONRenderer(sampleReport())
	result, err := renderer.Render(tmpFile)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	content, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.True(t, len(content) > 0)
	assert.Contains(t, string(content), "Sample Report")
}

func TestJSONRenderer_Render_MalformedFooter(t *testing.T) {
	// Footer is not a map, so timestamp injection should be skipped silently
	badReport := map[string]interface{}{
		"Title":  "Bad Report",
		"Footer": "not-a-map",
	}
	renderer := NewJSONRenderer(badReport).WithTimestamp(true)
	result, err := renderer.Render("")
	assert.NoError(t, err)

	assert.Contains(t, string(result), "Bad Report")
	assert.NotContains(t, string(result), "Generated at:")
}

func TestJSONRenderer_ImplementsRendererInterface(t *testing.T) {
	var _ interfaces.RendererInterface = NewJSONRenderer(sampleReport())
}
