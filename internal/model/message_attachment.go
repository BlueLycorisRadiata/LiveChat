package model

import "time"

type MessageAttachment struct {
	ID        int64     `json:"id" db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	FileURL   string    `json:"file_url" db:"file_url"`
	FileName  string    `json:"file_name" db:"file_name"`
	FileType  *string   `json:"file_type,omitempty" db:"file_type"`
	FileSize  *int64    `json:"file_size,omitempty" db:"file_size"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	Message *Message `json:"message,omitempty"`
}
