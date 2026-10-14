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

-- projects: cada aplicação que consome flags (ex.: agendaqui).
-- +goose StatementBegin
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  -- key é o identificador usado em URLs e SDKs: só minúsculas, números, "_" e "-".
  key TEXT NOT NULL UNIQUE CHECK (key ~ '^[a-z0-9][a-z0-9_-]*$'),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_projects
BEFORE UPDATE ON projects
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- environments: development, staging, production... de cada projeto.
-- +goose StatementBegin
CREATE TABLE environments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  -- ON DELETE CASCADE: apagar o projeto apaga os ambientes dele junto.
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name TEXT NOT NULL CHECK (name ~ '^[a-z0-9][a-z0-9_-]*$'),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  -- o mesmo nome pode existir em projetos diferentes, mas não duas vezes no mesmo.
  UNIQUE (project_id, name)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_environments
BEFORE UPDATE ON environments
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- feature_flags: a funcionalidade em si, sem estado. O estado fica em flag_configs.
-- +goose StatementBegin
CREATE TABLE feature_flags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  key TEXT NOT NULL CHECK (key ~ '^[a-z0-9][a-z0-9_-]*$'),
  description TEXT NOT NULL DEFAULT '',
  -- só boolean no MVP; novos tipos entram aqui depois (string, number, json).
  type TEXT NOT NULL DEFAULT 'boolean' CHECK (type IN ('boolean')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (project_id, key)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_feature_flags
BEFORE UPDATE ON feature_flags
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- flag_configs: o estado de uma flag em um ambiente (ligada em staging, desligada em production).
-- +goose StatementBegin
CREATE TABLE flag_configs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  flag_id UUID NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
  environment_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
  -- padrão conservador: flag nova nasce desligada.
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  -- regras de avaliação (V2). No MVP fica NULL.
  rules JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (flag_id, environment_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_timestamp_flag_configs
BEFORE UPDATE ON flag_configs
FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- O UNIQUE (flag_id, environment_id) já cria índice começando por flag_id.
-- Este cobre a busca no sentido contrário: "todas as configs de um ambiente".
-- +goose StatementBegin
CREATE INDEX idx_flag_configs_environment_id ON flag_configs(environment_id);
-- +goose StatementEnd

-- +goose Down
-- Ordem inversa da criação: quem tem FK sai antes de quem é referenciado.
-- (DROP TABLE também remove os triggers e índices da tabela.)
-- +goose StatementBegin
DROP TABLE IF EXISTS flag_configs;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS feature_flags;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS environments;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS projects;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS trigger_set_timestamp();
-- +goose StatementEnd
