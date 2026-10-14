package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"boliplate-go/internal/config"
	"boliplate-go/internal/db"
	"boliplate-go/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrations via goose. Comandos suportados (do próprio goose):
//
//	go run ./cmd/migrate up        # aplica todas as migrations pendentes
//	go run ./cmd/migrate down      # desfaz a última migration
//	go run ./cmd/migrate status    # mostra o estado de cada migration
//	go run ./cmd/migrate version   # mostra a versão atual
//	go run ./cmd/migrate redo      # desfaz e reaplica a última
//	go run ./cmd/migrate reset     # desfaz todas
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	command := "up"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	sqlDB, err := sql.Open("pgx", db.PostgresURL(cfg.DB))
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose: %v", err)
	}

	if err := goose.RunContext(context.Background(), command, sqlDB, ".", os.Args[2:]...); err != nil {
		log.Fatalf("goose %s: %v", command, err)
	}
}
