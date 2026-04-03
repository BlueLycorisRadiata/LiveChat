package handler

import (
	"LiveChat/internal/model"
	"LiveChat/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	convService service.ConversationService
	msgService  service.MessageService
}

func NewConversationHandler(convService service.ConversationService, msgService service.MessageService) *ConversationHandler {
	return &ConversationHandler{
		convService: convService,
		msgService:  msgService,
	}
}

type CreateConversationReq struct {
	Title          string  `json:"title"`
	ParticipantIDs []int64 `json:"participant_ids"`
	Type           string  `json:"type"`
}

type SendMessageReq struct {
	Content          string `json:"content"`
	Type             string `json:"type"`
	ReplyToMessageID *int64 `json:"reply_to_message_id"`
}

func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	userID := c.GetInt64("user_id")
	fmt.Printf("[DEBUG] CreateConversation userID: %d\n", userID)

	var req CreateConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convType := model.ConversationTypePrivate
	if req.Type == "group" {
		convType = model.ConversationTypeGroup
	}

	conv, err := h.convService.CreateConversation(c.Request.Context(), userID, req.Title, req.ParticipantIDs, convType)
	if err != nil {
		fmt.Printf("[DEBUG] CreateConversation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conv})
}

func (h *ConversationHandler) GetConversations(c *gin.Context) {
	userID := c.GetInt64("user_id")

	conversations, err := h.convService.GetUserConversations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conversations})
}

func (h *ConversationHandler) GetConversation(c *gin.Context) {
	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	conv, err := h.convService.GetConversation(c.Request.Context(), convIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": conv})
}

func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	userID := c.GetInt64("user_id")
	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	err = h.convService.DeleteConversation(c.Request.Context(), convIDInt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conversation deleted"})
}

func (h *ConversationHandler) LeaveConversation(c *gin.Context) {
	userID := c.GetInt64("user_id")
	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	err = h.convService.LeaveConversation(c.Request.Context(), convIDInt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "left conversation"})
}

func (h *ConversationHandler) GetMessages(c *gin.Context) {
	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	messages, err := h.msgService.GetMessages(c.Request.Context(), convIDInt, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

func (h *ConversationHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt64("user_id")

	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msgType := model.MessageTypeText
	if req.Type != "" {
		msgType = model.MessageType(req.Type)
	}

	msg, err := h.msgService.SendMessage(c.Request.Context(), convIDInt, userID, req.Content, msgType, req.ReplyToMessageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": msg})
}

func (h *ConversationHandler) DeleteMessage(c *gin.Context) {
	userID := c.GetInt64("user_id")

	msgID := c.Param("messageId")
	var msgIDInt int64
	_, err := fmt.Sscanf(msgID, "%d", &msgIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	err = h.msgService.DeleteMessage(c.Request.Context(), msgIDInt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}
