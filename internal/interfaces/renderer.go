// File: interfaces/renderer.go
package interfaces

// RendererInterface defines the contract for all renderers
// (e.g., PDF, JSON) that convert report structures into output formats.
type RendererInterface interface {

	// Render generates the output (PDF, JSON, etc.) and optionally writes
	// it to the given filename. If filename is empty, it returns the output as []byte.
	Render(filename string) ([]byte, error)

	// WithTimestamp enables or disables appending or embedding a timestamp
	// in the rendered output. Returns the renderer for chaining.
	WithTimestamp(enable bool) RendererInterface

	// SetTimestampFormat sets the timestamp layout using Go's time format layout rules.
	// Example: "2006-01-02 15:04:05" or time.RFC3339.
	SetTimestampFormat(format string) RendererInterface
}
