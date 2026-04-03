CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username VARCHAR NOT NULL UNIQUE,
                       email VARCHAR NOT NULL UNIQUE,
                       password VARCHAR NOT NULL
);

CREATE TYPE conversation_type AS ENUM ('private', 'group');
CREATE TYPE participant_role AS ENUM ('owner', 'member');
CREATE TYPE message_type AS ENUM ('text', 'image', 'file', 'system');

CREATE TABLE conversations (
                               id BIGSERIAL PRIMARY KEY,
                               type conversation_type NOT NULL,
                               title VARCHAR(255),
                               created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               deleted_at TIMESTAMPTZ
);

CREATE TABLE conversation_participants (
                                           id BIGSERIAL PRIMARY KEY,
                                           conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
                                           user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                           role participant_role NOT NULL DEFAULT 'member',
                                           joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                           left_at TIMESTAMPTZ,
                                           deleted_at TIMESTAMPTZ,
                                           last_read_message_id BIGINT,
                                           last_read_at TIMESTAMPTZ,
                                           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                           updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                           CONSTRAINT uq_conversation_participant UNIQUE (conversation_id, user_id)
);

CREATE TABLE messages (
                          id BIGSERIAL PRIMARY KEY,
                          conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
                          sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
                          content TEXT NOT NULL,
                          type message_type NOT NULL DEFAULT 'text',
                          reply_to_message_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          deleted_at TIMESTAMPTZ
);

ALTER TABLE conversation_participants
    ADD CONSTRAINT fk_last_read_message
        FOREIGN KEY (last_read_message_id)
            REFERENCES messages(id)
            ON DELETE SET NULL;

CREATE TABLE message_attachments (
                                     id BIGSERIAL PRIMARY KEY,
                                     message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
                                     file_url TEXT NOT NULL,
                                     file_name VARCHAR(255) NOT NULL,
                                     file_type VARCHAR(100),
                                     file_size BIGINT,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE message_reads (
                               id BIGSERIAL PRIMARY KEY,
                               message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
                               user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               CONSTRAINT uq_message_reads UNIQUE (message_id, user_id)
);

CREATE INDEX idx_messages_conversation_id
    ON messages(conversation_id, created_at);

CREATE INDEX idx_participants_user
    ON conversation_participants(user_id);

CREATE INDEX idx_participants_conversation
    ON conversation_participants(conversation_id);

CREATE INDEX idx_conversations_updated
    ON conversations(updated_at DESC);