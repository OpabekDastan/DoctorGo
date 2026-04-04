package postgres

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func AutoMigrate(db *sqlx.DB, dir string) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS app_migrations (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create app_migrations table: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("scan migrations: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", dir)
	}

	sort.Strings(files)

	for _, file := range files {
		name := filepath.Base(file)

		var alreadyApplied bool
		if err := db.Get(&alreadyApplied, `SELECT EXISTS(SELECT 1 FROM app_migrations WHERE name = $1)`, name); err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if alreadyApplied {
			log.Printf("migration skipped: %s", name)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}

		tx, err := db.Beginx()
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", name, err)
		}

		if _, err := tx.Exec(sql); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", name, err)
		}

		if _, err := tx.Exec(`INSERT INTO app_migrations (name) VALUES ($1)`, name); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}

		log.Printf("migration applied: %s", name)
	}

	return nil
}
