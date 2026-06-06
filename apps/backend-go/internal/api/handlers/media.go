package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"postpanda/backend-go/internal/api/middleware"
	"postpanda/backend-go/internal/config"
	"postpanda/backend-go/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaHandler struct {
	mediaSvc *services.MediaService
}

func NewMediaHandler(mediaSvc *services.MediaService) *MediaHandler {
	return &MediaHandler{mediaSvc: mediaSvc}
}

func (h *MediaHandler) List(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	mediaType := c.Query("type")

	media, total, err := h.mediaSvc.List(orgID, page, limit, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  media,
		"total": total,
	})
}

func (h *MediaHandler) Upload(c *gin.Context) {
	orgID := middleware.GetOrgID(c)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.App.MaxUploadSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No file provided"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	orgDir := filepath.Join(config.App.UploadDirectory, orgID)

	if err := os.MkdirAll(orgDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create upload directory"})
		return
	}

	destPath := filepath.Join(orgDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create file"})
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
		return
	}

	mediaType := detectMediaType(header.Header.Get("Content-Type"))
	relativePath := fmt.Sprintf("/uploads/%s/%s", orgID, filename)

	media, err := h.mediaSvc.Create(orgID, header.Filename, filename, relativePath, int(size), mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusCreated, media)
}

func (h *MediaHandler) Delete(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	if err := h.mediaSvc.Delete(orgID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted"})
}

func (h *MediaHandler) UpdateAlt(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	id := c.Param("id")

	var req struct {
		Alt string `json:"alt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.mediaSvc.UpdateAlt(orgID, id, req.Alt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alt text updated"})
}

func detectMediaType(contentType string) string {
	switch contentType {
	case "video/mp4", "video/webm", "video/mov":
		return "video"
	case "image/gif":
		return "gif"
	default:
		return "image"
	}
}

// Serve static files
func (h *MediaHandler) ServeFile(c *gin.Context) {
	orgID := c.Param("orgId")
	filename := c.Param("filename")

	// Sanitize path
	clean := filepath.Clean(filepath.Join(config.App.UploadDirectory, orgID, filename))
	expected := filepath.Join(config.App.UploadDirectory, orgID)
	if len(clean) < len(expected) {
		c.Status(http.StatusForbidden)
		return
	}

	http.ServeFile(c.Writer, c.Request, clean)
}

// Serve profile pictures
func (h *MediaHandler) SetProfilePicture(c *gin.Context) {
	orgID := middleware.GetOrgID(c)
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No file provided"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	orgDir := filepath.Join(config.App.UploadDirectory, orgID)
	os.MkdirAll(orgDir, 0755)

	destPath := filepath.Join(orgDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create file"})
		return
	}
	defer dst.Close()
	size, _ := io.Copy(dst, file)

	relativePath := fmt.Sprintf("/uploads/%s/%s", orgID, filename)
	media, err := h.mediaSvc.Create(orgID, header.Filename, filename, relativePath, int(size), "image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save media"})
		return
	}

	if err := h.mediaSvc.SetUserPicture(userID, media.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   media.ID,
		"path": relativePath,
		"time": time.Now(),
	})
}
