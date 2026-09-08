-- +goose Up

DROP TABLE IF EXISTS voice_messages;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS user_chats;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE TABLE chats (
    id UUID PRIMARY KEY,
    chat_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_chats (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, chat_id)
);

CREATE INDEX idx_user_chats_chat_id ON user_chats (chat_id);

CREATE TABLE messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    content TEXT,
    message_type VARCHAR(20) NOT NULL DEFAULT 'text',
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    reply_to_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT messages_message_type_check
        CHECK (message_type IN ('text', 'voice', 'video', 'file'))
);

CREATE INDEX idx_messages_chat_id_id ON messages (chat_id, id DESC);
CREATE INDEX idx_messages_chat_id_created_at ON messages (chat_id, created_at DESC);
CREATE INDEX idx_messages_reply_to_id ON messages (reply_to_id) WHERE reply_to_id IS NOT NULL;

CREATE TABLE voice_messages (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    object_key VARCHAR(512) NOT NULL,
    duration_seconds INT NOT NULL,
    waveform BYTEA,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT voice_messages_message_id_key UNIQUE (message_id),
    CONSTRAINT voice_messages_duration_seconds_check CHECK (duration_seconds >= 0)
);

CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    entity_name VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    payload BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT outbox_status_check
        CHECK (status IN ('pending', 'processed', 'failed'))
);

CREATE INDEX idx_outbox_entity_name_status_created_at
    ON outbox (entity_name, status, created_at);

-- +goose Down

DROP TABLE IF EXISTS voice_messages;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS user_chats;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
