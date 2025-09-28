-- +goose Up
CREATE TABLE comment_reactions (
    id SERIAL PRIMARY KEY,
    comment_id INT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(comment_id, user_id)
);

-- +goose Down
DROP TABLE comment_reactions;
