package repository

import (
	"LiveChat/internal/model"
	"context"
	"database/sql"
	"time"
)

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conv *model.Conversation) (*model.Conversation, error)
	GetConversationByID(ctx context.Context, id int64) (*model.Conversation, error)
	GetUserConversations(ctx context.Context, userID int64) ([]model.Conversation, error)
	UpdateConversation(ctx context.Context, id int64, userID int64, title string) (*model.Conversation, error)
	DeleteConversation(ctx context.Context, id int64, userID int64) error
	AddParticipant(ctx context.Context, convID int64, userID int64, role model.ParticipantRole) error
	RemoveParticipant(ctx context.Context, convID int64, userID int64) error
	GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error)
}

type conversationRepo struct {
	db DBTX
}

func NewConversationRepository(db DBTX) ConversationRepository {
	return &conversationRepo{db: db}
}

func (r *conversationRepo) CreateConversation(ctx context.Context, conv *model.Conversation) (*model.Conversation, error) {
	query := `
		INSERT INTO conversations (type, title, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		conv.Type, conv.Title, conv.CreatedBy, time.Now(), time.Now(),
	).Scan(&conv.ID)

	if err != nil {
		return nil, err
	}

	// Auto-insert default AI settings for AI conversations
	if conv.Type == model.ConversationTypeAI {
		aiQuery := `
			INSERT INTO conversation_ai_settings (conversation_id, model, temperature, max_tokens, system_prompt, created_at, updated_at)
			VALUES ($1, 'nvidia/nemotron-3-super-120b-a12b:free', 0.7, 2048, '', $2, $3)`
		_, err = r.db.ExecContext(ctx, aiQuery, conv.ID, time.Now(), time.Now())
		if err != nil {
			return nil, err
		}
	}

	return conv, nil
}

func (r *conversationRepo) GetConversationByID(ctx context.Context, id int64) (*model.Conversation, error) {
	query := `
		SELECT id, type, title, created_by, created_at, updated_at, deleted_at
		FROM conversations
		WHERE id = $1 AND deleted_at IS NULL`

	conv := &model.Conversation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID, &conv.Type, &conv.Title, &conv.CreatedBy,
		&conv.CreatedAt, &conv.UpdatedAt, &conv.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *conversationRepo) GetUserConversations(ctx context.Context, userID int64) ([]model.Conversation, error) {
	query := `
		SELECT c.id, c.type, c.title, c.created_by, c.created_at, c.updated_at, c.deleted_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1 AND cp.deleted_at IS NULL AND c.deleted_at IS NULL
		ORDER BY c.updated_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []model.Conversation
	for rows.Next() {
		var conv model.Conversation
		err := rows.Scan(
			&conv.ID, &conv.Type, &conv.Title, &conv.CreatedBy,
			&conv.CreatedAt, &conv.UpdatedAt, &conv.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *conversationRepo) DeleteConversation(ctx context.Context, id int64, userID int64) error {
	query := `
		UPDATE conversations
		SET deleted_at = NOW()
		WHERE id = $1 AND created_by = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *conversationRepo) UpdateConversation(ctx context.Context, id int64, userID int64, title string) (*model.Conversation, error) {
	query := `
		UPDATE conversations
		SET title = $1, updated_at = NOW()
		WHERE id = $2 AND created_by = $3 AND deleted_at IS NULL
		RETURNING id, type, title, created_by, created_at, updated_at, deleted_at`

	conv := &model.Conversation{}
	err := r.db.QueryRowContext(ctx, query, title, id, userID).Scan(
		&conv.ID, &conv.Type, &conv.Title, &conv.CreatedBy,
		&conv.CreatedAt, &conv.UpdatedAt, &conv.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *conversationRepo) AddParticipant(ctx context.Context, convID int64, userID int64, role model.ParticipantRole) error {
	query := `
		INSERT INTO conversation_participants (conversation_id, user_id, role, joined_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, convID, userID, role, time.Now(), time.Now(), time.Now())
	return err
}

func (r *conversationRepo) RemoveParticipant(ctx context.Context, convID int64, userID int64) error {
	query := `
		UPDATE conversation_participants
		SET left_at = NOW(), deleted_at = NOW()
		WHERE conversation_id = $1 AND user_id = $2 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, convID, userID)
	return err
}

func (r *conversationRepo) GetParticipants(ctx context.Context, convID int64) ([]model.ConversationParticipant, error) {
	query := `
		SELECT cp.id, cp.conversation_id, cp.user_id, cp.role, cp.joined_at, cp.left_at, 
		       cp.deleted_at, cp.last_read_message_id, cp.last_read_at, cp.created_at, cp.updated_at,
		       u.id, u.username, u.email
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.conversation_id = $1 AND cp.deleted_at IS NULL`

	rows, err := r.db.QueryContext(ctx, query, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.ConversationParticipant
	for rows.Next() {
		var p model.ConversationParticipant
		var u model.User
		err := rows.Scan(
			&p.ID, &p.ConversationID, &p.UserID, &p.Role, &p.JoinedAt, &p.LeftAt,
			&p.DeletedAt, &p.LastReadMessageID, &p.LastReadAt, &p.CreatedAt, &p.UpdatedAt,
			&u.ID, &u.Username, &u.Email,
		)
		if err != nil {
			return nil, err
		}
		p.User = &u
		participants = append(participants, p)
	}

	return participants, nil
}
