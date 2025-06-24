package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/interfaces/dto"
)

// UploadHandler handles file upload operations
type UploadHandler struct {
	uploadDir string
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler() *UploadHandler {
	uploadDir := "/app/uploads"
	// For local development, use relative path
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		uploadDir = "./uploads"
	}
	
	// Ensure upload directory exists
	os.MkdirAll(filepath.Join(uploadDir, "images", "articles"), 0755)
	
	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

// UploadImage handles image upload for articles
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルサイズが大きすぎます（最大10MB）"})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルが選択されていません"})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "サポートされていないファイル形式です（JPEG, PNG, WebPのみ）"})
		return
	}

	// Validate file size (5MB limit)
	if header.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ファイルサイズが大きすぎます（最大5MB）"})
		return
	}

	// Generate unique filename
	ext := getFileExtension(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Create date-based subdirectory
	now := time.Now()
	dateDir := fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	fullDir := filepath.Join(h.uploadDir, "images", "articles", dateDir)
	
	// Ensure directory exists
	err = os.MkdirAll(fullDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "アップロードディレクトリの作成に失敗しました"})
		return
	}

	// Save file
	filePath := filepath.Join(fullDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存に失敗しました"})
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルのコピーに失敗しました"})
		return
	}

	// Generate relative path for database storage
	relativePath := filepath.Join("images", "articles", dateDir, filename)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/") // Ensure forward slashes

	// Return response
	response := dto.ImageUploadResponse{
		ImagePath:     relativePath,
		ThumbnailPath: relativePath, // For now, same as original. TODO: Generate thumbnail
		OriginalName:  header.Filename,
		Size:          header.Size,
		ContentType:   contentType,
	}

	c.JSON(http.StatusOK, response)
}

// ServeImage serves uploaded images
func (h *UploadHandler) ServeImage(c *gin.Context) {
	imagePath := c.Param("path")
	
	// Security: prevent path traversal
	if strings.Contains(imagePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスです"})
		return
	}

	fullPath := filepath.Join(h.uploadDir, imagePath)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "画像が見つかりません"})
		return
	}

	// Serve file
	c.File(fullPath)
}

// DeleteImage deletes an uploaded image
func (h *UploadHandler) DeleteImage(c *gin.Context) {
	imagePath := c.Query("path")
	if imagePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "画像パスが指定されていません"})
		return
	}

	// Security: prevent path traversal
	if strings.Contains(imagePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なパスです"})
		return
	}

	fullPath := filepath.Join(h.uploadDir, imagePath)
	
	// Delete file
	err := os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "画像が見つかりません"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "画像の削除に失敗しました"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "画像を削除しました"})
}

// Helper functions

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg", 
		"image/png",
		"image/webp",
	}
	
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".jpg" // default extension
	}
	return strings.ToLower(ext)
}