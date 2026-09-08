package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedded embed.FS

func Up(db *sql.DB) error {
	goose.SetBaseFS(embedded)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			return nil
		}

		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
