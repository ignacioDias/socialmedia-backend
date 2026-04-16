package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type UploadHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

type uploadHandler struct {
}

func NewUploadHandler() UploadHandler {
	return &uploadHandler{}
}

const (
	uploadsDir   = "uploads"
	maxFileSize  = 10 << 20 // 10 MB
	imagesSubdir = "images"
)

var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func (h *uploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Ensure upload directory exists
	if err := os.MkdirAll(filepath.Join(uploadsDir, imagesSubdir), 0755); err != nil {
		http.Error(w, "Error preparing upload directory", http.StatusInternalServerError)
		return
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "Error parsing form or file too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read first 512 bytes to detect real MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	realMime := http.DetectContentType(buffer)
	ext, allowed := allowedTypes[realMime]
	if !allowed {
		http.Error(w, "Only JPEG and PNG are allowed", http.StatusUnsupportedMediaType)
		return
	}

	// Reset to beginning before saving
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadsDir, imagesSubdir, uuid.New().String()+ext)
	tempFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		// Clean up file on copy failure
		_ = os.Remove(filePath)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	WriteResponseWithEncoder(w, filePath, http.StatusCreated)
}
