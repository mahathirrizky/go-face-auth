package helper

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxUploadFileSize = 10 << 20 // 10 MB

var allowedUploadContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

// SaveUploadedFile saves an uploaded file to the specified destination.
// It returns the full path to the saved file and an error, if any.
// The destination path will be STORAGE_BASE_PATH/subDir/unique_filename.ext
func SaveUploadedFile(file *multipart.FileHeader, subDir string) (string, error) {
	if err := validateUploadedFile(file); err != nil {
		return "", err
	}

	storageBaseDir := os.Getenv("STORAGE_BASE_PATH")
	if storageBaseDir == "" {
		storageBaseDir = "/tmp/go_face_auth_data" // Fallback for development/testing
	}

	// Create the full destination directory
	destinationDir := filepath.Join(storageBaseDir, subDir)
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", destinationDir, err)
	}

	// Create a unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	uniqueFilename := uuid.New().String() + ext
	filePath := filepath.Join(destinationDir, uniqueFilename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file %s: %w", filePath, err)
	}
	defer dst.Close()

	// Copy the uploaded file content to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	return filePath, nil
}

func validateUploadedFile(file *multipart.FileHeader) error {
	if file == nil {
		return fmt.Errorf("file header is required")
	}
	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}
	if file.Size > maxUploadFileSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", maxUploadFileSize)
	}
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(file.Filename))
	}
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if _, ok := allowedUploadContentTypes[contentType]; !ok {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
	return nil
}
