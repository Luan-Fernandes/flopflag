# Arquitetura do projeto (Gin + GORM + Postgres + Migrations SQL via goose)

Este guia explica **passo a passo** como usar a arquitetura na prática, do zero até ter um CRUD completo funcionando (com exemplo real de `todos`).

## Visão geral

O projeto é uma API REST em camadas, com responsabilidades bem definidas:

- **`cmd/api`**: entrypoint da API HTTP.
- **`cmd/migrate`**: entrypoint de migrations (SQL) usando `goose` (migrations embarcadas via `embed.FS`).
- **`internal/config`**: leitura/validação de variáveis de ambiente.
- **`internal/db`**: helpers de conexão e URL do Postgres.
- **`internal/domain`**: entidades do domínio (structs usadas pelo GORM).
- **`internal/repository`**: acesso a dados (GORM), sem regras de negócio.
- **`internal/service`**: regras de negócio / casos de uso (validações, decisões).
- **`internal/controller`**: controllers HTTP (Gin): parse de request, chama service, devolve response.
- **`internal/httpserver`**: composição do servidor (router, middlewares, rotas).
- **`internal/httpx`**: envelope/handlers de resposta HTTP padronizados.
- **`migrations/`**: migrations SQL versionadas no formato goose (um arquivo por migration, com blocos `-- +goose Up` / `-- +goose Down`).

## Como clonar e configurar o ambiente

### Pré-requisitos

- **Go** instalado (projeto usa `go 1.23.0` no `go.mod`)
- **Postgres** disponível (local ou via Docker)

### Clonando o repositório

```bash
git clone <URL_DO_SEU_REPO>.git
cd boliplate-go
```

### Configurando variáveis de ambiente

Crie um `.env` a partir do exemplo:

```bash
cp .env.example .env
```

Ajuste pelo menos:

- **`DB_HOST`**: se você roda Postgres local, geralmente é `localhost` (no compose é `postgres`)
- **`DB_USER`**, **`DB_PASSWORD`**, **`DB_NAME`**
- **`JWT_SECRET`**: obrigatório para rodar a API

### Rodando o Postgres

Opção A) Subir Postgres via Docker (recomendado para dev):

```bash
docker compose up -d postgres
```

Opção B) Postgres local:

- Garanta um database criado e as credenciais batendo com seu `.env`.

## Como rodar o projeto (API + migrations)

### 1) Rodar migrations

```bash
go run ./cmd/migrate up
```

Comandos úteis (todos vêm do próprio goose):

- Ver estado de cada migration:

```bash
go run ./cmd/migrate status
```

- Ver versão atual:

```bash
go run ./cmd/migrate version
```

- Rollback da última migration:

```bash
go run ./cmd/migrate down
```

- Rollback total:

```bash
go run ./cmd/migrate reset
```

- Desfazer e reaplicar a última:

```bash
go run ./cmd/migrate redo
```

### 2) Rodar a API

```bash
go run ./cmd/api
```

Healthcheck:

```bash
curl -s http://localhost:8080/health
```

## Como criar novas tabelas (migrations)

**Regra de ouro:** não edite migrations antigas em projetos já compartilhados. Crie uma nova migration.

### Passo a passo

1) Crie **um** arquivo em `migrations/`, com o próximo número sequencial:

- `0000X_nome_descritivo.sql`

2) Dentro dele, escreva o `Up` (criar/alterar) e o `Down` (rollback) usando as anotações do goose.

3) Rode:

```bash
go run ./cmd/migrate up
```

Obs: como as migrations são embarcadas no binário via `embed.FS`, depois de criar/editar um `.sql` basta rodar `go run ./cmd/migrate up` (o `go run`/`go build` reembarca os arquivos automaticamente).

### Exemplo: adicionar uma coluna nova

`migrations/0000X_add_users_name.sql`:

```sql
-- +goose Up
ALTER TABLE users ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN name;
```

> Dica: statements que contêm `;` no meio (ex.: funções `plpgsql` com `$$ ... $$`) precisam ser
> envolvidos em `-- +goose StatementBegin` / `-- +goose StatementEnd`. Veja `migrations/00001_init.sql`.

## Como criar um novo módulo (boas práticas)

Um “módulo” aqui normalmente é um recurso REST (ex: `todos`) com:

- Entidade (domain)
- Migration SQL
- Repository (persistência)
- Service (regras)
- Controller (HTTP)
- Rotas (wiring)

### Checklist recomendado

- **`migrations/00000X_...`**: cria a tabela e constraints relevantes.
- **`internal/domain/<entidade>.go`**: struct com tags GORM e `TableName()`.
- **`internal/repository/<entidade>_repository.go`**:
  - Sem validações de negócio
  - Queries simples e legíveis
  - Evite N+1 quando houver relações (use `Preload` quando necessário)
- **`internal/service/<entidade>_service.go`**:
  - Validações (campos obrigatórios, ranges)
  - Regras (permissões, invariantes)
  - Erros de negócio viram `httpx.NewAppError(status, message, errs)`
- **`internal/controller/<entidade>_controller.go`**:
  - Parse de params/query/body
  - Chama service
  - Responde com `internal/httpx`
- **`internal/httpserver/routes.go`** + **`internal/httpserver/server.go`**:
  - Injeta repo/service/controller
  - Registra rotas

## Exemplo prático completo: módulo `todos`

Este projeto já inclui um módulo `todos` completo. Abaixo está o “caminho” que você repetiria para um módulo novo.

### 1) Migration: criar tabela `todos`

Arquivo:

- `migrations/00001_init.sql` (blocos `-- +goose Up` / `-- +goose Down`)

Obs: o `00001_init` cria `users` e `todos` (porque `todos` depende de `users`).

Aplicar:

```bash
go run ./cmd/migrate up
```

### 2) Domain: entidade `Todo`

Arquivo:

- `internal/domain/todo.go`

Responsabilidade:

- Define os campos persistidos
- Define `TableName()` para bater com o nome da tabela

### 3) Repository: acesso a dados

Arquivo:

- `internal/repository/todo_repository.go`

Responsabilidade:

- `Create`, `ListByUserID`, `FindByIDForUser`, `Update`, `Delete`

### 4) Service: regras/validações

Arquivo:

- `internal/service/todo_service.go`

Validações feitas:

- `title` obrigatório
- `Update` exige ao menos um campo
- `Get/Delete` retornam `404` quando não encontrado
- Todas as operações respeitam o `user_id` do token

### 5) Controller: CRUD HTTP (Gin)

Arquivo:

- `internal/controller/todo_controller.go`

Endpoints implementados (todos protegidos por JWT):

- `GET /todos`
- `POST /todos`
- `GET /todos/:id`
- `PATCH /todos/:id`
- `DELETE /todos/:id`

### 6) Rotas e wiring

Arquivos:

- `internal/httpserver/routes.go` (rotas)
- `internal/httpserver/server.go` (injeção: repo -> service -> controller)

## Como testar, rodar e validar o módulo `todos`

### 1) Subir Postgres + rodar migrations

Se estiver com Docker:

```bash
docker compose up -d postgres
go run ./cmd/migrate up
```

### 2) Rodar API

```bash
go run ./cmd/api
```

### 3) Obter um token (register/login)

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

Copie o `token` retornado e exporte:

```bash
export TOKEN="COLE_AQUI"
```

### 4) CRUD de todos (cURL)

Criar:

```bash
curl -s -X POST http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"primeira tarefa"}'
```

Listar:

```bash
curl -s http://localhost:8080/todos \
  -H "Authorization: Bearer $TOKEN"
```

Buscar por ID:

```bash
curl -s http://localhost:8080/todos/<id> \
  -H "Authorization: Bearer $TOKEN"
```

Atualizar:

```bash
curl -s -X PATCH http://localhost:8080/todos/<id> \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"done":true}'
```

Deletar:

```bash
curl -s -X DELETE http://localhost:8080/todos/<id> \
  -H "Authorization: Bearer $TOKEN" -i
```

### 5) Validação por testes automatizados

Hoje o repositório roda:

```bash
go test ./...
```

Esse comando valida compilação e testes existentes (ex: `internal/security/*_test.go`).

## Dicas para futuros módulos

- **Migrations são a fonte da verdade do schema**:
  - criar tabela/coluna/constraint = nova migration
  - rollback = `down.sql`
- **Controllers devem ser finos**:
  - parse/validação de entrada “estrutural” (JSON, query)
  - erro de negócio e regras ficam no service
- **Services devem ser o centro das regras**:
  - validações de domínio
  - retorno de `httpx.NewAppError(...)` para erros esperados
- **Repository não decide regra**:
  - só acessa dados
  - queries legíveis e sem duplicação
- **Resposta padronizada**:
  - use `internal/httpx` (OK/Fail helpers)

