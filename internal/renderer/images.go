package renderer

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

func mustSaveBase64ToTempFile(dataURI string) string {
	if dataURI == "" {
		return ""
	}
	path, err := saveBase64ImageToTempFile(dataURI)
	if err != nil {
		log.Printf("Failed to save base64 image: %v", err)
		return ""
	}
	return path
}

func saveBase64ImageToTempFile(dataURI string) (string, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URI")
	}

	prefix := parts[0]
	var ext string
	switch {
	case strings.Contains(prefix, "image/png"):
		ext = ".png"
	case strings.Contains(prefix, "image/jpeg"):
		ext = ".jpg"
	case strings.Contains(prefix, "image/svg+xml"):
		ext = ".svg"
	default:
		ext = ".img"
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	tmpfile, err := os.CreateTemp("", "img_*"+ext)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(decoded); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
