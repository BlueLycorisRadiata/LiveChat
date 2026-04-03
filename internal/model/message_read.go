package model

import "time"

type MessageRead struct {
	ID        int64     `json:"id" db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ReadAt    time.Time `json:"read_at" db:"read_at"`

	Message *Message `json:"message,omitempty"`
	User    *User    `json:"user,omitempty"`
}
