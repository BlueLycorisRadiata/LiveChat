package handler

import (
	"LiveChat/internal/ai"
	"LiveChat/internal/model"
	"LiveChat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	client *ai.Client
	aiSvc  service.AIService
}

func NewAIHandler(client *ai.Client, aiSvc service.AIService) *AIHandler {
	return &AIHandler{client: client, aiSvc: aiSvc}
}

func (h *AIHandler) ListModels(c *gin.Context) {
	models := ai.GetAllowedModelsList()
	c.JSON(http.StatusOK, gin.H{"data": models})
}

func (h *AIHandler) StreamMessage(c *gin.Context) {
	userID := c.GetInt64("user_id")

	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	err = h.aiSvc.StreamReply(c.Request.Context(), conversationID, userID, req.Content, c.Writer, flusher)
	if err != nil {
		c.Writer.Write([]byte("event: error\ndata: " + err.Error() + "\n\n"))
		flusher.Flush()
		return
	}

	c.Writer.Write([]byte("event: done\ndata: [DONE]\n\n"))
	flusher.Flush()
}

func (h *AIHandler) GetSettings(c *gin.Context) {
	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	settings, err := h.aiSvc.GetAISettings(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI settings not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": settings})
}

type UpdateAISettingsReq struct {
	Model        *string  `json:"model"`
	Temperature  *float64 `json:"temperature"`
	MaxTokens    *int     `json:"max_tokens"`
	SystemPrompt *string  `json:"system_prompt"`
}

func (h *AIHandler) UpdateSettings(c *gin.Context) {
	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req UpdateAISettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch current settings to merge partial update
	current, err := h.aiSvc.GetAISettings(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI settings not found"})
		return
	}

	if req.Model != nil {
		if !ai.IsModelAllowed(*req.Model) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "model not allowed"})
			return
		}
		current.Model = *req.Model
	}
	if req.Temperature != nil {
		current.Temperature = *req.Temperature
	}
	if req.MaxTokens != nil {
		current.MaxTokens = *req.MaxTokens
	}
	if req.SystemPrompt != nil {
		current.SystemPrompt = *req.SystemPrompt
	}

	err = h.aiSvc.UpdateAISettings(c.Request.Context(), conversationID, &model.ConversationAISettings{
		Model:        current.Model,
		Temperature:  current.Temperature,
		MaxTokens:    current.MaxTokens,
		SystemPrompt: current.SystemPrompt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "settings updated"})
}
