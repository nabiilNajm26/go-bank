package middleware

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	MaxFileSize = 5 * 1024 * 1024 // 5MB
)

var allowedImageTypes = []string{".png", ".jpg", ".jpeg"}

func FileUploadMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only apply to multipart form requests
		if !strings.Contains(c.Get("Content-Type"), "multipart/form-data") {
			return c.Next()
		}

		// Parse multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse multipart form",
			})
		}

		// Check all file fields
		for fieldName, files := range form.File {
			for _, fileHeader := range files {
				if err := validateFile(fileHeader); err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": fmt.Sprintf("Invalid file in field '%s': %v", fieldName, err),
					})
				}
			}
		}

		return c.Next()
	}
}

func validateFile(header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d (5MB)", header.Size, MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedImageType(ext) {
		return fmt.Errorf("invalid file type '%s'. Only PNG, JPG, JPEG allowed", ext)
	}

	// Check filename length
	if len(header.Filename) > 255 {
		return fmt.Errorf("filename too long (max 255 characters)")
	}

	// Check for dangerous filenames
	if strings.Contains(header.Filename, "..") || strings.Contains(header.Filename, "/") {
		return fmt.Errorf("invalid filename")
	}

	return nil
}

func isAllowedImageType(ext string) bool {
	for _, allowed := range allowedImageTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}