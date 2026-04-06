package model

import "time"

type ConversationAISettings struct {
	ID             int64     `json:"id" db:"id"`
	ConversationID int64     `json:"conversation_id" db:"conversation_id"`
	Model          string    `json:"model" db:"model"`
	Temperature    float64   `json:"temperature" db:"temperature"`
	MaxTokens      int       `json:"max_tokens" db:"max_tokens"`
	SystemPrompt   string    `json:"system_prompt" db:"system_prompt"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
