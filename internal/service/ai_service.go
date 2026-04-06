package service

import (
	"LiveChat/internal/ai"
	"LiveChat/internal/model"
	"LiveChat/internal/repository"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type AIService interface {
	StreamReply(ctx context.Context, convID int64, userID int64, userContent string, w http.ResponseWriter, flusher http.Flusher) error
	GetAISettings(ctx context.Context, convID int64) (*model.ConversationAISettings, error)
	UpdateAISettings(ctx context.Context, convID int64, settings *model.ConversationAISettings) error
}

type aiSvc struct {
	orClient     *ai.Client
	msgRepo      repository.MessageRepository
	settingsRepo repository.AISettingsRepository
	convRepo     repository.ConversationRepository
}

func NewAIService(
	orClient *ai.Client,
	msgRepo repository.MessageRepository,
	settingsRepo repository.AISettingsRepository,
	convRepo repository.ConversationRepository,
) AIService {
	return &aiSvc{
		orClient:     orClient,
		msgRepo:      msgRepo,
		settingsRepo: settingsRepo,
		convRepo:     convRepo,
	}
}

func (s *aiSvc) StreamReply(ctx context.Context, convID int64, userID int64, userContent string, w http.ResponseWriter, flusher http.Flusher) error {
	// 1. Validate conversation is type=ai
	conv, err := s.convRepo.GetConversationByID(ctx, convID)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}
	if conv.Type != model.ConversationTypeAI {
		return fmt.Errorf("conversation %d is not an AI conversation", convID)
	}

	// 2. Persist user message
	userRole := model.RoleUser
	_, err = s.msgRepo.CreateMessage(ctx, &model.Message{
		ConversationID: convID,
		SenderID:       userID,
		Content:        userContent,
		Type:           model.MessageTypeText,
		Role:           &userRole,
	})
	if err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// 3. Fetch AI settings
	settings, err := s.settingsRepo.GetByConversationID(ctx, convID)
	if err != nil {
		return fmt.Errorf("failed to get AI settings: %w", err)
	}

	// 4. Fetch last 30 messages for context
	history, err := s.settingsRepo.GetMessagesForAI(ctx, convID, 30)
	if err != nil {
		return fmt.Errorf("failed to get message history: %w", err)
	}

	// 5. Build messages array for OpenRouter
	var messages []ai.Message

	// Prepend system prompt if set
	if strings.TrimSpace(settings.SystemPrompt) != "" {
		messages = append(messages, ai.Message{
			Role:    "system",
			Content: settings.SystemPrompt,
		})
	}

	for _, m := range history {
		if m.Role == nil {
			continue
		}
		role := string(*m.Role)
		if role != "user" && role != "assistant" && role != "system" {
			continue
		}
		if strings.TrimSpace(m.Content) == "" {
			continue
		}
		messages = append(messages, ai.Message{
			Role:    role,
			Content: m.Content,
		})
	}

	// 6. Stream response via SSE
	fullAnswer := ""

	err = s.orClient.StreamChat(ctx, ai.ChatRequest{
		Model:       settings.Model,
		Messages:    messages,
		Temperature: settings.Temperature,
		MaxTokens:   settings.MaxTokens,
	}, func(token string) {
		fullAnswer += token
		w.Write([]byte("data: " + token + "\n\n"))
		flusher.Flush()
	}, func(usage *ai.Usage) {
		// optional: log usage
	})

	if err != nil {
		return fmt.Errorf("stream failed: %w", err)
	}

	// 7. Persist assistant response in DB
	assistantRole := model.RoleAssistant
	_, err = s.msgRepo.CreateMessage(ctx, &model.Message{
		ConversationID: convID,
		SenderID:       0,
		Content:        fullAnswer,
		Type:           model.MessageTypeText,
		Role:           &assistantRole,
	})
	if err != nil {
		return fmt.Errorf("failed to save assistant message: %w", err)
	}

	return nil
}

func (s *aiSvc) GetAISettings(ctx context.Context, convID int64) (*model.ConversationAISettings, error) {
	return s.settingsRepo.GetByConversationID(ctx, convID)
}

func (s *aiSvc) UpdateAISettings(ctx context.Context, convID int64, settings *model.ConversationAISettings) error {
	settings.ConversationID = convID
	return s.settingsRepo.Upsert(ctx, settings)
}
