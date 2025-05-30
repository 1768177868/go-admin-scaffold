package handlers

import (
	"app/internal/core/storage"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler 处理文件上传
type UploadHandler struct {
	storage storage.Storage
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(storage storage.Storage) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

// UploadRequest 上传请求
type UploadRequest struct {
	Type string `form:"type" binding:"required,oneof=avatar image file"` // 文件类型：avatar-头像，image-图片，file-其他文件
}

// UploadResponse 上传响应
type UploadResponse struct {
	URL  string `json:"url"`  // 文件访问URL
	Path string `json:"path"` // 文件存储路径
}

// Upload 处理文件上传
func (h *UploadHandler) Upload(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("failed to get file: %v", err),
		})
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	if !isAllowedFileType(req.Type, ext) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "file type not allowed",
		})
		return
	}

	// 验证文件大小
	if !isAllowedFileSize(req.Type, file.Size) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "file size exceeds limit",
		})
		return
	}

	// 构建存储路径
	path := buildStoragePath(req.Type)

	// 上传文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("failed to open file: %v", err),
		})
		return
	}
	defer src.Close()

	filePath := filepath.Join(path, file.Filename)
	url, err := h.storage.Put(filePath, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("failed to upload file: %v", err),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": UploadResponse{
			URL:  url,
			Path: filePath,
		},
	})
}

// 检查文件类型是否允许
func isAllowedFileType(fileType, ext string) bool {
	switch fileType {
	case "avatar":
		return isImageExt(ext)
	case "image":
		return isImageExt(ext)
	case "file":
		return true
	default:
		return false
	}
}

// 检查是否是图片扩展名
func isImageExt(ext string) bool {
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return allowedExts[ext]
}

// 检查文件大小是否允许
func isAllowedFileSize(fileType string, size int64) bool {
	switch fileType {
	case "avatar":
		return size <= 2*1024*1024 // 2MB
	case "image":
		return size <= 5*1024*1024 // 5MB
	case "file":
		return size <= 10*1024*1024 // 10MB
	default:
		return false
	}
}

// 构建存储路径
func buildStoragePath(fileType string) string {
	now := time.Now()
	return filepath.Join(fileType, fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day()))
}
