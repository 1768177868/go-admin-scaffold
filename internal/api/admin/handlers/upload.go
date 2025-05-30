package handlers

import (
	"app/internal/core/storage"
	"app/pkg/response"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file uploads
type UploadHandler struct {
	storage storage.Storage
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage storage.Storage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

// Upload handles file uploads
// @Summary Upload file
// @Description Upload a file to storage (local or S3)
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param type formData string true "File type (image, document, video, audio)"
// @Success 200 {object} map[string]string
// @Router /upload/file [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	// 获取文件类型
	fileType := c.PostForm("type")
	if fileType == "" {
		response.ParamError(c, "type is required")
		return
	}

	// 验证文件类型
	allowedTypes := map[string][]string{
		"image":    {".jpg", ".jpeg", ".png", ".gif", ".webp"},
		"document": {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"},
		"video":    {".mp4", ".avi", ".mov", ".wmv"},
		"audio":    {".mp3", ".wav", ".ogg", ".m4a"},
	}

	extensions, ok := allowedTypes[fileType]
	if !ok {
		response.ParamError(c, fmt.Sprintf("invalid file type. allowed types: %s", strings.Join(getKeys(allowedTypes), ", ")))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.ParamError(c, "file is required")
		return
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !contains(extensions, ext) {
		response.ParamError(c, fmt.Sprintf("invalid file extension. allowed extensions: %s", strings.Join(extensions, ", ")))
		return
	}

	// 验证文件大小
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		response.ParamError(c, fmt.Sprintf("file too large. maximum size: %dMB", maxSize/1024/1024))
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.ServerError(c)
		return
	}
	defer src.Close()

	// 生成存储路径
	now := time.Now()
	path := fmt.Sprintf("%s/%d/%02d/%02d/%s", fileType, now.Year(), now.Month(), now.Day(), file.Filename)

	// 上传文件
	url, err := h.storage.Put(path, src)
	if err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, gin.H{
		"url":  url,
		"path": path,
		"name": file.Filename,
		"size": file.Size,
		"type": fileType,
	})
}

// getKeys returns the keys of a map
func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// contains checks if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
