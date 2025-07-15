// main_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_JSONOutput(t *testing.T) {
	tmpDir := t.TempDir()

	// Prepare test files
	inputFile := filepath.Join(tmpDir, "input.json")
	tmplFile := filepath.Join(tmpDir, "template.html")
	outputFile := filepath.Join(tmpDir, "report.json")

	err := os.WriteFile(inputFile, []byte(`{"Title": "Test"}`), 0644)
	if err != nil {
		return
	}
	err = os.WriteFile(tmplFile, []byte(`{{.Title}}`), 0644)
	if err != nil {
		return
	}

	args := []string{
		"--input", inputFile,
		"--template", tmplFile,
		"--format", "json",
		"--output", outputFile,
	}

	err = Run(args)
	assert.NoError(t, err)

	data, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Test")
}

func TestRun_MissingRequiredArgs(t *testing.T) {
	err := Run([]string{"--input", "a.json"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required")
}

func TestRun_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create dummy input and template files
	inputPath := filepath.Join(tmpDir, "input.json")
	templatePath := filepath.Join(tmpDir, "template.html")

	err := os.WriteFile(inputPath, []byte(`{"key":"value"}`), 0644)
	if err != nil {
		return
	}
	err = os.WriteFile(templatePath, []byte(`{{.key}}`), 0644)
	if err != nil {
		return
	}

	args := []string{
		"--input", inputPath,
		"--template", templatePath,
		"--format", "txt", // <-- invalid format
	}

	err = Run(args)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid format")
}
