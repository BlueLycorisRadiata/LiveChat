package service

import (
	"LiveChat/internal/model"
	"LiveChat/internal/repository"
	"context"
	"fmt"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, userID int64, title string, participantIDs []int64, convType model.ConversationType) (*model.Conversation, error)
	GetConversation(ctx context.Context, convID int64) (*model.Conversation, error)
	GetUserConversations(ctx context.Context, userID int64) ([]model.Conversation, error)
	DeleteConversation(ctx context.Context, convID int64, userID int64) error
	LeaveConversation(ctx context.Context, convID int64, userID int64) error
	GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error)
}

type conversationSvc struct {
	convRepo repository.ConversationRepository
	msgRepo  repository.MessageRepository
}

func NewConversationService(convRepo repository.ConversationRepository, msgRepo repository.MessageRepository) ConversationService {
	return &conversationSvc{
		convRepo: convRepo,
		msgRepo:  msgRepo,
	}
}

func (s *conversationSvc) CreateConversation(ctx context.Context, userID int64, title string, participantIDs []int64, convType model.ConversationType) (*model.Conversation, error) {
	fmt.Printf("[DEBUG] Service CreateConversation userID=%d, title=%s, type=%s\n", userID, title, convType)

	conv := &model.Conversation{
		Type:      convType,
		CreatedBy: userID,
	}

	if title != "" {
		conv.Title = &title
	} else if convType == model.ConversationTypeAI {
		defaultTitle := "AI Chat"
		conv.Title = &defaultTitle
	}

	createdConv, err := s.convRepo.CreateConversation(ctx, conv)
	if err != nil {
		fmt.Printf("[DEBUG] CreateConversation repo error: %v\n", err)
		return nil, err
	}

	// Add creator as owner
	err = s.convRepo.AddParticipant(ctx, createdConv.ID, userID, model.ParticipantRoleOwner)
	if err != nil {
		fmt.Printf("[DEBUG] AddParticipant error: %v\n", err)
		return nil, err
	}

	// Add other participants (skip for AI conversations — only the creator is needed)
	if convType != model.ConversationTypeAI {
		for _, pID := range participantIDs {
			if pID != userID {
				err = s.convRepo.AddParticipant(ctx, createdConv.ID, pID, model.ParticipantRoleMember)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return createdConv, nil
}

func (s *conversationSvc) GetConversation(ctx context.Context, convID int64) (*model.Conversation, error) {
	return s.convRepo.GetConversationByID(ctx, convID)
}

func (s *conversationSvc) GetUserConversations(ctx context.Context, userID int64) ([]model.Conversation, error) {
	return s.convRepo.GetUserConversations(ctx, userID)
}

func (s *conversationSvc) DeleteConversation(ctx context.Context, convID int64, userID int64) error {
	return s.convRepo.DeleteConversation(ctx, convID, userID)
}

func (s *conversationSvc) LeaveConversation(ctx context.Context, convID int64, userID int64) error {
	return s.convRepo.RemoveParticipant(ctx, convID, userID)
}

func (s *conversationSvc) GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error) {
	return s.convRepo.GetParticipants(ctx, convID)
}
