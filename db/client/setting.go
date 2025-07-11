package client

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manthan307/corebase/db/schema"
)

type SettingRepo interface {
	Create(ctx context.Context, param schema.SettingParams) error
	Get(ctx context.Context, key string) (*schema.SettingRow, error)
	Update(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]schema.SettingRow, error)
	GetBunch(ctx context.Context, key []string) ([]*schema.SettingRow, error)
}

type setting struct {
	db *pgxpool.Pool
}

func NewSettingRepo(db *pgxpool.Pool) SettingRepo {
	return &setting{db: db}
}

// Create inserts or upserts a new setting
func (s *setting) Create(ctx context.Context, param schema.SettingParams) error {
	query := `
		INSERT INTO private.settings ("key", "value")
		VALUES ($1, $2)
		ON CONFLICT ("key") DO UPDATE SET
			"value" = EXCLUDED."value",
			updated_at = now();`

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Rollback on error

	_, err = tx.Exec(ctx, query, param.Key, param.Value)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get retrieves a setting by key
func (s *setting) Get(ctx context.Context, key string) (*schema.SettingRow, error) {
	query := `SELECT "key", "value" FROM private.settings WHERE "key" = $1 AND deleted_at IS NULL LIMIT 1`

	row := s.db.QueryRow(ctx, query, key)

	var setting schema.SettingRow
	if err := row.Scan(&setting.Key, &setting.Value); err != nil {
		return nil, err
	}

	return &setting, nil
}

// Update modifies the value of an existing setting
func (s *setting) Update(ctx context.Context, key string, value string) error {
	query := `UPDATE private.settings SET "value" = $1, updated_at = now() WHERE "key" = $2 AND deleted_at IS NULL`

	ct, err := s.db.Exec(ctx, query, value, key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("setting not found or already deleted")
	}
	return nil
}

// Delete soft-deletes a setting
func (s *setting) Delete(ctx context.Context, key string) error {
	query := `UPDATE private.settings SET deleted_at = $1, updated_at = now() WHERE "key" = $2 AND deleted_at IS NULL`

	ct, err := s.db.Exec(ctx, query, time.Now(), key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("setting not found or already deleted")
	}
	return nil
}

// List returns all non-deleted settings
func (s *setting) List(ctx context.Context) ([]schema.SettingRow, error) {
	query := `SELECT "key", "value" FROM private.settings WHERE deleted_at IS NULL ORDER BY "key"`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []schema.SettingRow
	for rows.Next() {
		var s schema.SettingRow
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}

	return settings, nil
}

// GetBunch retrieves multiple settings by a list of keys
func (s *setting) GetBunch(ctx context.Context, keys []string) ([]*schema.SettingRow, error) {
	if len(keys) == 0 {
		return []*schema.SettingRow{}, nil // Return empty slice, not nil
	}

	query := `
		SELECT "key", "value"
		FROM private.settings
		WHERE "key" = ANY($1::text[])
		  AND deleted_at IS NULL
		ORDER BY "key";
	`

	rows, err := s.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*schema.SettingRow
	for rows.Next() {
		var row schema.SettingRow
		if err := rows.Scan(&row.Key, &row.Value); err != nil {
			return nil, err
		}
		result = append(result, &row)
	}

	return result, nil
}
