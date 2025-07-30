// File: renderer/fonts.go
package pkg

import (
	"fmt"
	"github.com/signintech/gopdf"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// findFontPaths returns the full paths to font files located in the assets directory.
//
// It infers the base directory relative to this file's location and looks for:
//   - LiberationSans-Regular.ttf
//   - LiberationSans-Bold.ttf
//   - LiberationSans-Italic.ttf
//
// If a font file is not found, it logs a warning and returns an empty string for that font.
func findFontPaths() (FontPaths, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return FontPaths{}, fmt.Errorf("cannot determine caller location")
	}

	baseDir := filepath.Dir(currentFile)
	assetsDir := filepath.Join(baseDir, "..", "..", "assets")

	check := func(filename string) string {
		path := filepath.Join(assetsDir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
		log.Printf("Font not found: %s", path)
		return ""
	}

	return FontPaths{
		Regular: check("LiberationSans-Regular.ttf"),
		Bold:    check("LiberationSans-Bold.ttf"),
		Italic:  check("LiberationSans-Italic.ttf"),
	}, nil
}

// SetFont applies a dynamic font style (regular, bold, italic, or bold-italic) to the PDF context.
//
// It assumes the font family is named "Arial" (registered via AddTTFFontWithOption).
// This function chooses the correct style string based on the italic and bold flags.
func SetFont(pdf *gopdf.GoPdf, italic, bold bool) {
	style := ""
	if bold {
		style += "B"
	}
	if italic {
		style += "I"
	}
	_ = pdf.SetFont("Arial", style, 12)
}
