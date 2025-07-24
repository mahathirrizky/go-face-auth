package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"strings"

	"github.com/google/uuid"
)

// SaveUploadedFile saves an uploaded file to the specified destination.
// It returns the full path to the saved file and an error, if any.
// The destination path will be STORAGE_BASE_PATH/subDir/unique_filename.ext
func SaveUploadedFile(file *multipart.FileHeader, subDir string) (string, error) {
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
