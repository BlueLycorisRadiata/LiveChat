package model

import "time"

type Conversation struct {
	ID        int64            `json:"id" db:"id"`
	Type      ConversationType `json:"type" db:"type"`
	Title     *string          `json:"title,omitempty" db:"title"`
	CreatedBy int64            `json:"created_by" db:"created_by"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time       `json:"deleted_at,omitempty" db:"deleted_at"`

	Participants []ConversationParticipant `json:"participants,omitempty"`
	Messages     []Message                 `json:"messages,omitempty"`
}
