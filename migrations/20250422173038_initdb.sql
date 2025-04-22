-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    users TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE messages (
    from_user TEXT NOT NULL CHECK (length(from_user) <= 40),
    text TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;
DROP TABLE messages;
-- +goose StatementEnd
