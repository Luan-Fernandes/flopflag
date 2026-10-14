## Boilerplate Go (Gin + GORM + Postgres)

API REST simples e “profissional”, sem complexidade desnecessária:

- **Gin** como framework HTTP
- **Postgres** via `docker compose`
- **GORM** como ORM
- **Migrations SQL** com `goose` (embarcadas no binário via `embed.FS`)
- **JWT** (login/register) + **roles** (`user`/`admin`)
- **Envelope de resposta** padronizado
- **CRUD exemplo**: `todos`

### Documentação (passo a passo)

Veja o guia completo em `docs/ARQUITETURA.md`.

### Como rodar

Opcionalmente, crie um `.env` a partir do exemplo (o `docker compose` já tem defaults):

```bash
cp .env.example .env
```

Suba tudo (banco + migrations + API):

```bash
docker compose --profile app up --build
```

Rodar local (sem container), assumindo Postgres já rodando:

```bash
go run ./cmd/migrate up
go run ./cmd/api
```

Subir apenas o banco via Docker:

```bash
docker compose up -d postgres
```

### Endpoints

- **Health**
  - `GET /health`
- **Auth**
  - `POST /auth/register`
  - `POST /auth/login`
- **Autenticado (Bearer token)**
  - `GET /me`
  - `GET /todos`
  - `POST /todos`
  - `GET /todos/:id`
  - `PATCH /todos/:id`
  - `DELETE /todos/:id`
- **Admin**
  - `GET /admin/users?limit=50&offset=0`

Obs: **o primeiro usuário registrado** vira `admin` (bootstrap simples).

### Exemplo rápido (cURL)

Registrar:

```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"12345678"}'
```

Login:

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"12345678"}'
```

Criar todo (substitua `$TOKEN`):

```bash
curl -s -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"primeira tarefa"}'
```

### Estrutura

- `cmd/api`: entrypoint HTTP
- `internal/config`: config via env
- `internal/db`: conexão GORM/Postgres
- `internal/domain`: entidades
- `internal/repository`: acesso a dados
- `internal/service`: regras de negócio
- `internal/controller`: controllers (camada HTTP)
- `internal/httpserver`: rotas, middleware, envelope
- `migrations`: migrations SQL (formato goose) embarcadas via `embed.FS` e rodadas via `cmd/migrate`

# boilerplate-api-go
