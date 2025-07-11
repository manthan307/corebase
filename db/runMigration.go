package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manthan307/corebase/db/schema"
	"go.uber.org/zap"
)

var (
	query = []string{
		`CREATE SCHEMA IF NOT EXISTS private;`,
		`CREATE SCHEMA IF NOT EXISTS public;`,
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE OR REPLACE FUNCTION private.set_updated_at() RETURNS TRIGGER AS $$
			BEGIN
 			   NEW.updated_at = now();
  			   RETURN NEW;
			END;
		$$ LANGUAGE plpgsql;`,
	}
	indexQuery = []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_settings_key ON private.settings (key);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS admin_username_key ON private.admin (username);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS admin_email_key ON private.admin (email);`,
	}
)

func RunMigration(db *pgxpool.Pool, log *zap.Logger) {
	query = append(query, schema.SettingQuery...)
	for _, q := range query {
		_, err := db.Exec(context.Background(), q)
		if err != nil {
			log.Fatal("Error running migration", zap.Error(err))
		}
	}

	for _, q := range indexQuery {
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Fatal("Error starting transaction for index migration", zap.Error(err))
		}
		_, err = tx.Exec(context.Background(), q)
		if err != nil {
			tx.Rollback(context.Background())
			log.Fatal("Error running index migration", zap.Error(err))
		}
		tx.Commit(context.Background())
	}

	log.Info("✅ Database migrations ran successfully")

}
