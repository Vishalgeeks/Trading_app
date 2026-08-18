package design

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SaveUploadedFile(r *http.Request, fieldName, uploadDir string) (string, error) {
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to get uploaded file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".bin"
	}

	nameBuf := make([]byte, 16)
	if _, err := rand.Read(nameBuf); err != nil {
		return "", fmt.Errorf("failed to generate filename: %w", err)
	}
	filename := hex.EncodeToString(nameBuf) + ext
	path := filepath.Join(uploadDir, filename)

	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return "/uploads/" + filename, nil
}
