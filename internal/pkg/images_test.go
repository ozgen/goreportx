package pkg

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMustSaveBase64ToTempFile_ValidImage(t *testing.T) {
	base64PNG := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQImWNgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="

	path := mustSaveBase64ToTempFile(base64PNG)
	assert.NotEmpty(t, path)

	info, err := os.Stat(path)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())

	_ = os.Remove(path) // cleanup
}

func TestMustSaveBase64ToTempFile_EmptyInput(t *testing.T) {
	path := mustSaveBase64ToTempFile("")
	assert.Equal(t, "", path)
}

func TestSaveBase64ImageToTempFile_ValidJPG(t *testing.T) {
	base64JPG := "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxISEhUTEhIVFhUVFRUVFRUVFRUVFRUWFxUWFhUYHSggGBolGxUXITEhJSkrLi4uFx8zODMtNygtLisBCgoKDg0OGhAQGy0lICYtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLf/AABEIAMgAyAMBIgACEQEDEQH/xAAcAAACAwEBAQEAAAAAAAAAAAAEBQADBgECBwj/xABAEAABAwIDBQUEBwUGBQUAAAABAAIRAyEEEjFBBVFhBhMicYGRoQYjQrHB0fAVUnKS4SNCUvEWI2OiwvEHFjZTgqKy/8QAGgEBAQEAAwEAAAAAAAAAAAAAAAECBAUGA//EADsRAQACAQMCAwUGBAQHAAAAAAABAgMABBEFEiExBkFRcYETIjKhscHR8FJSwSNCYnLh8YIVM0OCFiRjk//aAAwDAQACEQMRAD8A"

	path, err := saveBase64ImageToTempFile(base64JPG)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, ".jpg"))

	_ = os.Remove(path)
}

func TestSaveBase64ImageToTempFile_InvalidBase64(t *testing.T) {
	invalid := "data:image/png;base64,NOT-REAL-BASE64"

	path, err := saveBase64ImageToTempFile(invalid)
	assert.Error(t, err)
	assert.Equal(t, "", path)
}

func TestSaveBase64ImageToTempFile_InvalidDataURI(t *testing.T) {
	invalid := "not-a-valid-uri"

	path, err := saveBase64ImageToTempFile(invalid)
	assert.Error(t, err)
	assert.Equal(t, "", path)
}
