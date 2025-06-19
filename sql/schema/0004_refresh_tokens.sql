-- +goose Up
CREATE TABLE refresh_tokens (
    token PRIMARY KEY TEXT,
    created_at TIMESTAMP NOT NUll,
    updated_at TIMESTAMP NOT NULL,
    user NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at timestamp
);

-- +goose Down
DROP TABLE refresh_tokens;