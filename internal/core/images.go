// File: renderer/image.go
package core

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

// mustSaveBase64ToTempFile attempts to decode and save a base64-encoded image (data URI)
// to a temporary file. It returns the path to the saved image.
// If an error occurs, it logs the issue and returns an empty string.
//
// This function is useful for inline image support (e.g., <img src="data:image/png;base64,...">).
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

// saveBase64ImageToTempFile parses a data URI, decodes the base64 image,
// determines the appropriate file extension (e.g., .png, .jpg), and writes
// the binary data to a temporary file.
//
// Returns the path to the created temporary file or an error if the decoding/writing fails.
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
