package pdf

import (
	"bytes"
	"github.com/ozgen/goreportx/pkg/core"
	"github.com/ozgen/goreportx/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"html/template"
	"os"
	"testing"
)

type MockRenderer struct {
	core.Renderer
	RenderFunc func(content string) (*bytes.Buffer, error)
}

func (m *MockRenderer) RenderHTMLLikeToBuffer(content string) (*bytes.Buffer, error) {
	if m.RenderFunc != nil {
		return m.RenderFunc(content)
	}
	return bytes.NewBufferString("PDF_BYTES:" + content), nil
}

type MockFactory struct {
	core.RendererFactory
	Renderer *core.Renderer
	Err      error
}

func (f *MockFactory) Build() (*core.Renderer, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Renderer, nil
}

func setupTestRenderer(t *testing.T, tmplStr string, factory *core.RendererFactory) *PDFRenderer {
	tmpl := template.Must(template.New("test").Parse(tmplStr))
	return NewPDFRenderer(
		map[string]interface{}{"Title": "My Report"},
		tmpl,
		factory,
	).(*PDFRenderer)
}

func TestPDFRenderer_Render_Basic(t *testing.T) {
	mockRenderer := &MockRenderer{
		RenderFunc: func(content string) (*bytes.Buffer, error) {
			return bytes.NewBufferString("Generated PDF: " + content), nil
		},
	}
	mockFactory := &MockFactory{Renderer: &mockRenderer.Renderer}

	pdfRenderer := setupTestRenderer(t, `Title: {{.Title}}`, &mockFactory.RendererFactory)
	_, err := pdfRenderer.Render("")

	assert.NoError(t, err)
}

func TestPDFRenderer_Render_WithTimestamp(t *testing.T) {
	mockRenderer := &MockRenderer{
		RenderFunc: func(content string) (*bytes.Buffer, error) {
			return bytes.NewBufferString("PDF with timestamp"), nil
		},
	}
	mockFactory := &MockFactory{Renderer: &mockRenderer.Renderer}

	pdfRenderer := setupTestRenderer(t, `Hi`, &mockFactory.RendererFactory)
	pdfRenderer.WithTimestamp(true)
	_, err := pdfRenderer.Render("")

	assert.NoError(t, err)
}

func TestPDFRenderer_Render_WithCustomTimestampFormat(t *testing.T) {
	customFormat := "02-Jan-2006"
	mockRenderer := &MockRenderer{
		RenderFunc: func(content string) (*bytes.Buffer, error) {
			return bytes.NewBufferString("Custom timestamp OK"), nil
		},
	}
	mockFactory := &MockFactory{Renderer: &mockRenderer.Renderer}

	pdfRenderer := setupTestRenderer(t, `Hi`, &mockFactory.RendererFactory)
	pdfRenderer.WithTimestamp(true).SetTimestampFormat(customFormat)
	_, err := pdfRenderer.Render("")

	assert.NoError(t, err)
}

func TestPDFRenderer_Render_WriteToFile(t *testing.T) {
	mockRenderer := &MockRenderer{
		RenderFunc: func(content string) (*bytes.Buffer, error) {
			return bytes.NewBufferString("SavedPDF"), nil
		},
	}
	mockFactory := &MockFactory{Renderer: &mockRenderer.Renderer}

	tmpFile := "test_report.pdf"
	defer os.Remove(tmpFile)

	pdfRenderer := setupTestRenderer(t, `test`, &mockFactory.RendererFactory)
	_, err := pdfRenderer.Render(tmpFile)

	assert.NoError(t, err)

	_, err = os.ReadFile(tmpFile)
	assert.NoError(t, err)
}

func TestPDFRenderer_ImplementsRendererInterface(t *testing.T) {
	var _ interfaces.RendererInterface = NewPDFRenderer(nil, template.New("x"), &core.RendererFactory{})
}
