-- +goose Up
ALTER TABLE chirp
ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
ALTER TABLE chirp
ALTER COLUMN user_id DROP NOT NULL;