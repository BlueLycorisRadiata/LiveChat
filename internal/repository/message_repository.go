package repository

import (
	"LiveChat/internal/model"
	"context"
	"database/sql"
	"time"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *model.Message) (*model.Message, error)
	GetMessagesByConversation(ctx context.Context, convID int64, limit, offset int) ([]model.Message, error)
	GetMessageByID(ctx context.Context, id int64) (*model.Message, error)
	DeleteMessage(ctx context.Context, id int64, userID int64) error
}

type messageRepo struct {
	db DBTX
}

func NewMessageRepository(db DBTX) MessageRepository {
	return &messageRepo{db: db}
}

func (r *messageRepo) CreateMessage(ctx context.Context, msg *model.Message) (*model.Message, error) {
	query := `
		INSERT INTO messages (conversation_id, sender_id, content, type, reply_to_message_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		msg.ConversationID, msg.SenderID, msg.Content, msg.Type,
		msg.ReplyToMessageID, time.Now(), time.Now(),
	).Scan(&msg.ID)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *messageRepo) GetMessagesByConversation(ctx context.Context, convID int64, limit, offset int) ([]model.Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.type, 
		       m.reply_to_message_id, m.created_at, m.updated_at, m.deleted_at,
		       u.id, u.username, u.email
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1 AND m.deleted_at IS NULL
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, convID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var m model.Message
		var u model.User
		err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.Type,
			&m.ReplyToMessageID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&u.ID, &u.Username, &u.Email,
		)
		if err != nil {
			return nil, err
		}
		m.Sender = &u
		messages = append(messages, m)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepo) GetMessageByID(ctx context.Context, id int64) (*model.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, content, type, 
		       reply_to_message_id, created_at, updated_at, deleted_at
		FROM messages
		WHERE id = $1 AND deleted_at IS NULL`

	msg := &model.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.Type,
		&msg.ReplyToMessageID, &msg.CreatedAt, &msg.UpdatedAt, &msg.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (r *messageRepo) DeleteMessage(ctx context.Context, id int64, userID int64) error {
	query := `
		UPDATE messages
		SET deleted_at = NOW()
		WHERE id = $1 AND sender_id = $2 AND deleted_at IS NULL`

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
