# flopflag

Serviço de feature flags: gerencie funcionalidades por projeto e ambiente e consulte o estado delas pela API, sem precisar de um novo deploy.

> Em construção: MVP.

## Stack

- **Go** + **Gin** (HTTP)
- **PostgreSQL** + **GORM**
- **goose** para migrations SQL (embarcadas no binário via `embed.FS`)
- **JWT** para autenticar o painel administrativo
- **Docker Compose** para desenvolvimento

## Como rodar

```bash
cp .env.example .env
docker compose --profile app up --build   # banco + migrations + API
```

Sem container, com o Postgres já rodando:

```bash
docker compose up -d postgres
go run ./cmd/migrate up
go run ./cmd/api
```

## Endpoints

- `GET /health`
- **Auth:** `POST /auth/register`, `POST /auth/login`
- **Autenticado (Bearer token):** `GET /me`
- **Admin:** `GET /admin/users?limit=50&offset=0`

O primeiro usuário registrado vira `admin`.
