package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spectrogram-app/internal/service"
)

type Handler struct {
	svc *service.AnalyzerService
	r   *gin.Engine
}

func NewHandler(svc *service.AnalyzerService) *Handler {
	h := &Handler{
		svc: svc,
		r:   gin.Default(),
	}
	h.setupRoutes()
	return h
}

func (h *Handler) setupRoutes() {
	h.r.POST("/analyze", h.analyzeImage)
	h.r.POST("/analyze-batch", h.analyzeBatch)
}

func (h *Handler) Run(addr string) error {
	return h.r.Run(addr)
}

func (h *Handler) analyzeImage(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.AnalyzeImage(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) analyzeBatch(c *gin.Context) {
	var req struct {
		URLs []string `json:"urls" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Реализовать пакетный анализ
	c.JSON(http.StatusOK, gin.H{"status": "not implemented"})
}
