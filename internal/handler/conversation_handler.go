package handler

import (
	"LiveChat/config"
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
	Model          string  `json:"model"`
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
	} else if req.Type == "ai" {
		convType = model.ConversationTypeAI
	}

	model := req.Model
	if model != "" && !config.IsModelSupported(model) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported AI model"})
		return
	}

	conv, err := h.convService.CreateConversation(c.Request.Context(), userID, req.Title, req.ParticipantIDs, convType, model)
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

type UpdateConversationReq struct {
	Title string `json:"title"`
}

func (h *ConversationHandler) UpdateConversation(c *gin.Context) {
	userID := c.GetInt64("user_id")
	convID := c.Param("id")
	fmt.Printf("[DEBUG] UpdateConversation: userID=%d, convID=%s\n", userID, convID)
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req UpdateConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("[DEBUG] UpdateConversation: title=%s\n", req.Title)

	conv, err := h.convService.UpdateConversation(c.Request.Context(), convIDInt, userID, req.Title)
	if err != nil {
		fmt.Printf("[DEBUG] UpdateConversation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

type AddMemberReq struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (h *ConversationHandler) AddMember(c *gin.Context) {
	requesterID := c.GetInt64("user_id")

	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	var req AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.convService.AddMember(c.Request.Context(), convIDInt, requesterID, req.UserID)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()
		if errMsg == "cannot add members to a private conversation" ||
			errMsg == "cannot add members to an AI conversation" ||
			errMsg == "user is already a member of this conversation" {
			status = http.StatusBadRequest
		} else if errMsg == "only admins or owners can add members" ||
			errMsg == "you are not a participant of this conversation" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

func (h *ConversationHandler) RemoveMember(c *gin.Context) {
	requesterID := c.GetInt64("user_id")

	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	targetUserID := c.Param("userId")
	var targetUserIDInt int64
	_, err = fmt.Sscanf(targetUserID, "%d", &targetUserIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.convService.RemoveMember(c.Request.Context(), convIDInt, requesterID, targetUserIDInt)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := err.Error()
		if errMsg == "cannot remove members from a private conversation" ||
			errMsg == "cannot remove members from an AI conversation" ||
			errMsg == "cannot remove the owner of the conversation" {
			status = http.StatusBadRequest
		} else if errMsg == "only admins or owners can remove members" ||
			errMsg == "only the owner can remove an admin" ||
			errMsg == "you are not a participant of this conversation" {
			status = http.StatusForbidden
		} else if errMsg == "target user is not a participant of this conversation" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

func (h *ConversationHandler) GetMembers(c *gin.Context) {
	convID := c.Param("id")
	var convIDInt int64
	_, err := fmt.Sscanf(convID, "%d", &convIDInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}

	members, err := h.convService.GetMembers(c.Request.Context(), convIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": members})
}
