package schema

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AdminParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AdminUpdateParams struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

var AdminQuery = []string{
	// Create enum type for roles
	`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
			CREATE TYPE role AS ENUM ('admin', 'superadmin');
		END IF;
	END$$;`,

	// Create the admin table
	`CREATE TABLE IF NOT EXISTS private.admin (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		role role NOT NULL DEFAULT 'admin',
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		deleted_at TIMESTAMPTZ
	);`,

	// Add trigger for updated_at on update
	`DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_trigger WHERE tgname = 'admin_set_updated_at'
		) THEN
			CREATE TRIGGER admin_set_updated_at
			BEFORE UPDATE ON private.admin
			FOR EACH ROW
			EXECUTE FUNCTION private.set_updated_at();
		END IF;
	END$$;`,
}
