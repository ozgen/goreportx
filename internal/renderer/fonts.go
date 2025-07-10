package renderer

import (
	"fmt"
	"github.com/signintech/gopdf"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type FontPaths struct {
	Regular string
	Bold    string
	Italic  string
}

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
