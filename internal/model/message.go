package model

import "time"

type Message struct {
	ID               int64       `json:"id" db:"id"`
	ConversationID   int64       `json:"conversation_id" db:"conversation_id"`
	SenderID         int64       `json:"sender_id" db:"sender_id"`
	Content          string      `json:"content" db:"content"`
	Type             MessageType `json:"type" db:"type"`
	ReplyToMessageID *int64      `json:"reply_to_message_id,omitempty" db:"reply_to_message_id"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time  `json:"deleted_at,omitempty" db:"deleted_at"`

	Sender       *User               `json:"sender,omitempty"`
	Conversation *Conversation       `json:"conversation,omitempty"`
	Attachments  []MessageAttachment `json:"attachments,omitempty"`
}
