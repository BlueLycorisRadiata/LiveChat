package service

import (
	"LiveChat/internal/model"
	"LiveChat/internal/repository"
	"context"
)

type MessageService interface {
	SendMessage(ctx context.Context, convID int64, senderID int64, content string, msgType model.MessageType, replyTo *int64) (*model.Message, error)
	GetMessages(ctx context.Context, convID int64, limit, offset int) ([]model.Message, error)
	DeleteMessage(ctx context.Context, msgID int64, userID int64) error
}

type messageSvc struct {
	msgRepo  repository.MessageRepository
	convRepo repository.ConversationRepository
}

func NewMessageService(msgRepo repository.MessageRepository, convRepo repository.ConversationRepository) MessageService {
	return &messageSvc{
		msgRepo:  msgRepo,
		convRepo: convRepo,
	}
}

func (s *messageSvc) SendMessage(ctx context.Context, convID int64, senderID int64, content string, msgType model.MessageType, replyTo *int64) (*model.Message, error) {
	msg := &model.Message{
		ConversationID:   convID,
		SenderID:         senderID,
		Content:          content,
		Type:             msgType,
		ReplyToMessageID: replyTo,
	}

	return s.msgRepo.CreateMessage(ctx, msg)
}

func (s *messageSvc) GetMessages(ctx context.Context, convID int64, limit, offset int) ([]model.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.msgRepo.GetMessagesByConversation(ctx, convID, limit, offset)
}

func (s *messageSvc) DeleteMessage(ctx context.Context, msgID int64, userID int64) error {
	return s.msgRepo.DeleteMessage(ctx, msgID, userID)
}
