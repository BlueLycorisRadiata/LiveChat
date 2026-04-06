package service

import (
	"LiveChat/config"
	"LiveChat/internal/model"
	"LiveChat/internal/repository"
	"context"
	"fmt"
)

type ConversationService interface {
	CreateConversation(ctx context.Context, userID int64, title string, participantIDs []int64, convType model.ConversationType, aiModel string) (*model.Conversation, error)
	GetConversation(ctx context.Context, convID int64) (*model.Conversation, error)
	GetUserConversations(ctx context.Context, userID int64) ([]model.Conversation, error)
	UpdateConversation(ctx context.Context, convID int64, userID int64, title string) (*model.Conversation, error)
	DeleteConversation(ctx context.Context, convID int64, userID int64) error
	LeaveConversation(ctx context.Context, convID int64, userID int64) error
	GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error)
}

type conversationSvc struct {
	convRepo       repository.ConversationRepository
	msgRepo        repository.MessageRepository
	aiSettingsRepo repository.AISettingsRepository
}

func NewConversationService(convRepo repository.ConversationRepository, msgRepo repository.MessageRepository, aiSettingsRepo repository.AISettingsRepository) ConversationService {
	return &conversationSvc{
		convRepo:       convRepo,
		msgRepo:        msgRepo,
		aiSettingsRepo: aiSettingsRepo,
	}
}

func (s *conversationSvc) CreateConversation(ctx context.Context, userID int64, title string, participantIDs []int64, convType model.ConversationType, aiModel string) (*model.Conversation, error) {
	fmt.Printf("[DEBUG] Service CreateConversation userID=%d, title=%s, type=%s, model=%s\n", userID, title, convType, aiModel)

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

	err = s.convRepo.AddParticipant(ctx, createdConv.ID, userID, model.ParticipantRoleOwner)
	if err != nil {
		fmt.Printf("[DEBUG] AddParticipant error: %v\n", err)
		return nil, err
	}

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

	if convType == model.ConversationTypeAI {
		selectedModel := aiModel
		if selectedModel == "" {
			selectedModel = config.GetDefaultAIModel()
		}
		defaultSettings := &model.ConversationAISettings{
			ConversationID: createdConv.ID,
			Model:          selectedModel,
			Temperature:    0.7,
			MaxTokens:      2048,
			SystemPrompt:   "",
		}
		err = s.aiSettingsRepo.Upsert(ctx, defaultSettings)
		if err != nil {
			fmt.Printf("[DEBUG] CreateAIsettings error: %v\n", err)
			return nil, fmt.Errorf("failed to create AI settings: %w", err)
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

func (s *conversationSvc) UpdateConversation(ctx context.Context, convID int64, userID int64, title string) (*model.Conversation, error) {
	return s.convRepo.UpdateConversation(ctx, convID, userID, title)
}

func (s *conversationSvc) LeaveConversation(ctx context.Context, convID int64, userID int64) error {
	return s.convRepo.RemoveParticipant(ctx, convID, userID)
}

func (s *conversationSvc) GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error) {
	return s.convRepo.GetParticipants(ctx, convID)
}
