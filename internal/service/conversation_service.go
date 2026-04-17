package service

import (
	"LiveChat/config"
	"LiveChat/internal/model"
	"LiveChat/internal/repository"
	"context"
	"errors"
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
	AddMember(ctx context.Context, convID int64, requesterID int64, targetUserID int64) error
	RemoveMember(ctx context.Context, convID int64, requesterID int64, targetUserID int64) error
	UpdateMemberRole(ctx context.Context, convID int64, requesterID int64, targetUserID int64, role model.ParticipantRole) error
	GetMembers(ctx context.Context, convID int64) ([]model.ConversationParticipant, error)
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

	if convType == model.ConversationTypePrivate {
		// Private chats must have exactly one other participant
		filtered := make([]int64, 0)
		for _, pID := range participantIDs {
			if pID != userID {
				filtered = append(filtered, pID)
			}
		}
		if len(filtered) != 1 {
			return nil, errors.New("private conversations must have exactly one other participant")
		}
		err = s.convRepo.AddParticipant(ctx, createdConv.ID, filtered[0], model.ParticipantRoleMember)
		if err != nil {
			return nil, err
		}
	} else if convType != model.ConversationTypeAI {
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

func (s *conversationSvc) AddMember(ctx context.Context, convID int64, requesterID int64, targetUserID int64) error {
	conv, err := s.convRepo.GetConversationByID(ctx, convID)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	if conv.Type == model.ConversationTypePrivate {
		return errors.New("cannot add members to a private conversation")
	}
	if conv.Type == model.ConversationTypeAI {
		return errors.New("cannot add members to an AI conversation")
	}

	alreadyIn, err := s.convRepo.IsParticipant(ctx, convID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to check participant: %w", err)
	}
	if alreadyIn {
		return errors.New("user is already a member of this conversation")
	}

	return s.convRepo.AddParticipant(ctx, convID, targetUserID, model.ParticipantRoleMember)
}

func (s *conversationSvc) RemoveMember(ctx context.Context, convID int64, requesterID int64, targetUserID int64) error {
	conv, err := s.convRepo.GetConversationByID(ctx, convID)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	if conv.Type == model.ConversationTypePrivate {
		return errors.New("cannot remove members from a private conversation")
	}
	if conv.Type == model.ConversationTypeAI {
		return errors.New("cannot remove members from an AI conversation")
	}

	return s.convRepo.RemoveParticipant(ctx, convID, targetUserID)
}

func (s *conversationSvc) UpdateMemberRole(ctx context.Context, convID int64, requesterID int64, targetUserID int64, role model.ParticipantRole) error {
	// Get conversation to check type
	conv, err := s.convRepo.GetConversationByID(ctx, convID)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}

	if conv.Type != model.ConversationTypeGroup {
		return errors.New("roles can only be updated in group conversations")
	}

	// Validate the role
	if !role.IsValid() {
		return errors.New("invalid role")
	}
	if role == model.ParticipantRoleOwner {
		return errors.New("cannot assign owner role")
	}

	// Only owner can change roles
	requesterRole, err := s.convRepo.GetParticipantRole(ctx, convID, requesterID)
	if err != nil {
		return errors.New("you are not a participant of this conversation")
	}
	if requesterRole != model.ParticipantRoleOwner {
		return errors.New("only the owner can change member roles")
	}

	// Cannot change own role
	if requesterID == targetUserID {
		return errors.New("cannot change your own role")
	}

	// Verify target is a participant
	exists, err := s.convRepo.IsParticipant(ctx, convID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to check participant: %w", err)
	}
	if !exists {
		return errors.New("target user is not a participant of this conversation")
	}

	return s.convRepo.UpdateParticipantRole(ctx, convID, targetUserID, role)
}

func (s *conversationSvc) GetMembers(ctx context.Context, convID int64) ([]model.ConversationParticipant, error) {
	return s.convRepo.GetParticipants(ctx, convID)
}
