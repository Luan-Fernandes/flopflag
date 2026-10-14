package migrations

import "embed"

// FS embarca todos os arquivos .sql de migration no binário,
// para que `cmd/migrate` não dependa da pasta em disco em runtime.
//
//go:embed *.sql
var FS embed.FS
