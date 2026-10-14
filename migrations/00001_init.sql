-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS todos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  done BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_timestamp_todos ON todos;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_todos
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todos;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS trigger_set_timestamp();
-- +goose StatementEnd
