package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manthan307/corebase/db/schema"
)

type AdminsRepo interface {
	Create(ctx context.Context, param schema.AdminParams) (schema.Admin, error)
	Get(ctx context.Context, id string) (schema.Admin, error)
	GetByEmail(ctx context.Context, email string) (schema.Admin, error)
	GetByUsername(ctx context.Context, username string) (schema.Admin, error)
	Update(ctx context.Context, id string, param schema.AdminUpdateParams) (schema.Admin, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]schema.Admin, error)
	TotalCount(ctx context.Context) (int, error)
	Exists(ctx context.Context, email, username string) (bool, string, error)
}

type admin struct {
	Db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) AdminsRepo {
	return &admin{Db: db}
}

func (a *admin) Create(ctx context.Context, param schema.AdminParams) (schema.Admin, error) {
	query := `
		INSERT INTO private.admin (username, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, role, created_at, updated_at;
	`

	tx, err := a.Db.Begin(ctx)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, param.Username, param.Email, param.Password, param.Role)

	var adminRow schema.Admin
	err = row.Scan(
		&adminRow.ID,
		&adminRow.Username,
		&adminRow.Email,
		&adminRow.Role,
		&adminRow.CreatedAt,
		&adminRow.UpdatedAt,
	)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("scan admin create: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return schema.Admin{}, fmt.Errorf("commit tx: %w", err)
	}

	return adminRow, nil
}

func (a *admin) Get(ctx context.Context, id string) (schema.Admin, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM private.admin
		WHERE id = $1 AND deleted_at IS NULL;
	`
	row := a.Db.QueryRow(ctx, query, id)

	var adminRow schema.Admin
	err := row.Scan(
		&adminRow.ID,
		&adminRow.Username,
		&adminRow.Email,
		&adminRow.Role,
		&adminRow.CreatedAt,
		&adminRow.UpdatedAt,
	)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("get admin by id: %w", err)
	}

	return adminRow, nil
}

func (a *admin) GetByEmail(ctx context.Context, email string) (schema.Admin, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM private.admin
		WHERE email = $1 AND deleted_at IS NULL;
	`
	row := a.Db.QueryRow(ctx, query, email)

	var admin schema.Admin
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.Role,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("get admin by email: %w", err)
	}

	return admin, nil
}

func (a *admin) GetByUsername(ctx context.Context, username string) (schema.Admin, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM private.admin
		WHERE username = $1 AND deleted_at IS NULL;
	`
	row := a.Db.QueryRow(ctx, query, username)

	var admin schema.Admin
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.Role,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("get admin by username: %w", err)
	}

	return admin, nil
}

func (a *admin) Update(ctx context.Context, id string, param schema.AdminUpdateParams) (schema.Admin, error) {
	query := `
		UPDATE private.admin
		SET
			username = COALESCE($2, username),
			email = COALESCE($3, email),
			password = COALESCE($4, password),
			role = COALESCE($5, role),
			updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, username, email, role, created_at, updated_at;
	`
	row := a.Db.QueryRow(ctx, query, id, param.Username, param.Email, param.Password, param.Role)

	var adminRow schema.Admin
	err := row.Scan(
		&adminRow.ID,
		&adminRow.Username,
		&adminRow.Email,
		&adminRow.Role,
		&adminRow.CreatedAt,
		&adminRow.UpdatedAt,
	)
	if err != nil {
		return schema.Admin{}, fmt.Errorf("update admin: %w", err)
	}

	return adminRow, nil
}

func (a *admin) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE private.admin
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL;
	`
	_, err := a.Db.Exec(ctx, query, id)
	return fmt.Errorf("delete admin: %w", err)
}

func (a *admin) List(ctx context.Context, limit, offset int) ([]schema.Admin, error) {
	query := `
		SELECT id, username, email, role, created_at, updated_at
		FROM private.admin
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`
	rows, err := a.Db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list admins: %w", err)
	}
	defer rows.Close()

	var admins []schema.Admin
	for rows.Next() {
		var adminRow schema.Admin
		err := rows.Scan(
			&adminRow.ID,
			&adminRow.Username,
			&adminRow.Email,
			&adminRow.Role,
			&adminRow.CreatedAt,
			&adminRow.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan admin row: %w", err)
		}
		admins = append(admins, adminRow)
	}

	return admins, nil
}

func (a *admin) TotalCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM private.admin
		WHERE deleted_at IS NULL;
	`
	var count int
	err := a.Db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count admins: %w", err)
	}
	return count, nil
}

// Exists checks if a username or email already exists.
func (a *admin) Exists(ctx context.Context, email, username string) (bool, string, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM private.admin
			WHERE (email = $1 OR username = $2)
			AND deleted_at IS NULL
		);
	`
	err := a.Db.QueryRow(ctx, query, email, username).Scan(&exists)
	if err != nil {
		return false, "", err
	}

	// Determine which one exists if needed (optional)
	if exists {
		query = `
			SELECT email, username
			FROM private.admin
			WHERE (email = $1 OR username = $2)
			AND deleted_at IS NULL
			LIMIT 1;
		`
		var e, u string
		err := a.Db.QueryRow(ctx, query, email, username).Scan(&e, &u)
		if err != nil {
			return true, "", err
		}

		if e == email {
			return true, "email", nil
		}
		if u == username {
			return true, "username", nil
		}
		return true, "both", nil
	}

	return false, "", nil
}
