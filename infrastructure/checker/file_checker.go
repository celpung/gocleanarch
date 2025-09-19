package checker

import (
	"log"
	"mime/multipart"
	"net/http"
)

// Helper function to check file type
func FileChecker(fileHeader *multipart.FileHeader) (bool, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 bytes of the file
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false, err
	}

	// Determine the content type of the file
	contentType := http.DetectContentType(buffer)

	// Check if the content type is an image or PDF
	if contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif" || contentType == "application/pdf" {
		return true, nil
	}

	// Optionally, log the invalid content type
	log.Printf("Invalid file type detected: %s", contentType)

	return false, nil
}
