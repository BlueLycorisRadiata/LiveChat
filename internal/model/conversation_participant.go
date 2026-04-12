package model

import "time"

type ConversationParticipant struct {
	ID                int64      `json:"id" db:"id"`
	ConversationID    int64      `json:"conversation_id" db:"conversation_id"`
	UserID            int64      `json:"user_id" db:"user_id"`
	JoinedAt          time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt            *time.Time `json:"left_at,omitempty" db:"left_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	LastReadMessageID *int64     `json:"last_read_message_id,omitempty" db:"last_read_message_id"`
	LastReadAt        *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`

	User         *User         `json:"user,omitempty"`
	Conversation *Conversation `json:"conversation,omitempty"`
}
