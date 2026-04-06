-- Add role column to messages for AI conversations (user | assistant | system)
ALTER TABLE messages ADD COLUMN IF NOT EXISTS role VARCHAR(20);

-- Conversation AI settings table (per-room AI model configuration)
CREATE TABLE IF NOT EXISTS conversation_ai_settings (
    id                BIGSERIAL PRIMARY KEY,
    conversation_id   BIGINT NOT NULL UNIQUE REFERENCES conversations(id) ON DELETE CASCADE,
    model             VARCHAR(255) NOT NULL DEFAULT 'nvidia/nemotron-3-super-120b-a12b:free',
    temperature       NUMERIC(4,2) NOT NULL DEFAULT 0.7,
    max_tokens        INT NOT NULL DEFAULT 2048,
    system_prompt     TEXT NOT NULL DEFAULT '',
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_conversation_ai_settings_conversation_id ON conversation_ai_settings(conversation_id);
