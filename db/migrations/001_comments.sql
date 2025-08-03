-- +goose Up
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 
    CONSTRAINT unique_user UNIQUE (user_id)
);

-- +goose Down
DROP TABLE comments;