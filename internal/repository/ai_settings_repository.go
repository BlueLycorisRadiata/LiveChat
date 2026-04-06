package repository

import (
	"LiveChat/internal/model"
	"context"
	"time"
)

type AISettingsRepository interface {
	GetByConversationID(ctx context.Context, convID int64) (*model.ConversationAISettings, error)
	Upsert(ctx context.Context, settings *model.ConversationAISettings) error
	GetMessagesForAI(ctx context.Context, convID int64, limit int) ([]model.Message, error)
}

type aiSettingsRepo struct {
	db DBTX
}

func NewAISettingsRepository(db DBTX) AISettingsRepository {
	return &aiSettingsRepo{db: db}
}

func (r *aiSettingsRepo) GetByConversationID(ctx context.Context, convID int64) (*model.ConversationAISettings, error) {
	query := `
		SELECT id, conversation_id, model, temperature, max_tokens, system_prompt, created_at, updated_at
		FROM conversation_ai_settings
		WHERE conversation_id = $1`

	s := &model.ConversationAISettings{}
	err := r.db.QueryRowContext(ctx, query, convID).Scan(
		&s.ID, &s.ConversationID, &s.Model, &s.Temperature,
		&s.MaxTokens, &s.SystemPrompt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *aiSettingsRepo) Upsert(ctx context.Context, settings *model.ConversationAISettings) error {
	query := `
		INSERT INTO conversation_ai_settings (conversation_id, model, temperature, max_tokens, system_prompt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (conversation_id)
		DO UPDATE SET model = $2, temperature = $3, max_tokens = $4, system_prompt = $5, updated_at = $7`

	_, err := r.db.ExecContext(ctx, query,
		settings.ConversationID, settings.Model, settings.Temperature,
		settings.MaxTokens, settings.SystemPrompt, time.Now(), time.Now(),
	)
	return err
}

func (r *aiSettingsRepo) GetMessagesForAI(ctx context.Context, convID int64, limit int) ([]model.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, type, role, created_at
		FROM messages
		WHERE conversation_id = $1 AND deleted_at IS NULL AND role IN ('user', 'assistant', 'system')
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, convID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var m model.Message
		err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.Type, &m.Role, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	// Reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
