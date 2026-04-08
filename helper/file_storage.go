package helper

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const maxUploadFileSize = 10 << 20 // 10 MB

var allowedUploadContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
}

var s3Client *s3.Client

func init() {
	// Initialize S3 Client during package init
	// It uses OS Env variables to configure Custom Endpoint for MinIO / Rust-S3 
	endpoint := os.Getenv("S3_ENDPOINT") // e.g. https://s3.yoursite.com
	region := os.Getenv("S3_REGION")     // e.g. us-east-1
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")

	if endpoint != "" && accessKey != "" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)

		if err != nil {
			fmt.Printf("Warning: fallback S3 init failed: %v\n", err)
		} else {
			// Path style is often required for self-hosted S3 clones
			// BaseEndpoint is the modern way to specify a custom endpoint in aws-sdk-go-v2
			s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(endpoint)
				o.UsePathStyle = true
			})
		}
	}
}

// SaveUploadedFile saves an uploaded file to the specified S3 destination.
// It returns the full public S3 URL and an error, if any.
// The destination path will be S3_BUCKET/subDir/unique_filename.ext
func SaveUploadedFile(file *multipart.FileHeader, subDir string) (string, error) {
	if err := validateUploadedFile(file); err != nil {
		return "", err
	}

	if s3Client == nil {
		return "", fmt.Errorf("S3 client is not initialized, check your S3_* environment variables")
	}

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return "", fmt.Errorf("S3_BUCKET environment variable is missing")
	}

	// Create a unique filename and prefix (object key)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	uniqueFilename := uuid.New().String() + ext
	
	// Ensure no leading slashes in S3 keys
	objectKey := filepath.Join(subDir, uniqueFilename)
	objectKey = strings.TrimPrefix(objectKey, "/") 

	// Open the uploaded memory/temp file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}

	// Upload to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(contentType),
		// ACL public-read is optional depending on bucket policy, usually bucket policy is preferred
		// ACL: types.ObjectCannedACLPublicRead, 
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Construct the final public URL
	endpoint := strings.TrimSuffix(os.Getenv("S3_ENDPOINT"), "/")
	fileUrl := fmt.Sprintf("%s/%s/%s", endpoint, bucket, objectKey)

	return fileUrl, nil
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
