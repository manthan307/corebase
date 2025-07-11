package schema

import (
	"time"

	"github.com/google/uuid"
)

type Settings struct {
	ID        uuid.UUID  `json:"id"`
	Key       string     `json:"key"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type SettingParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SettingRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var SettingQuery = []string{
	`CREATE TABLE IF NOT EXISTS private.settings (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		"key" TEXT UNIQUE NOT NULL,
		"value" TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		deleted_at TIMESTAMPTZ DEFAULT NULL
	);`,
	`DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_trigger 
		WHERE tgname = 'settings_set_updated_at'
	) THEN
		CREATE TRIGGER settings_set_updated_at
		BEFORE UPDATE ON private.settings
		FOR EACH ROW
		EXECUTE FUNCTION private.set_updated_at();
	END IF;
END$$;
`,
	`
	INSERT INTO private.settings ("key", "value")
VALUES
  ('auth.ban_threshold', '20'),
  ('auth.ban_duration', '24h'),
  ('auth.lockout_duration', '600s'),
  ('auth.token_expiry', '3600'),
  ('auth.max_login_attempts', '5'),
  ('app.name', 'corebase'),
  ('notifications.enabled', 'true'),
  ('rate_limit', '10'),
  ('server.read_timeout', '15'),
  ('server.write_timeout', '15'),
  ('server.idle_timeout', '60'),
  ('cors.enabled', 'true'),
  ('cors.allow_methods', '["GET", "POST", "PUT", "DELETE"]'),
  ('cors.allow_headers', '["Content-Type", "Authorization"]'),
  ('cors.allow_credentials', 'true'),
  ('uploads.max_file_size', '5242880'),
  ('uploads.allowed_types', '["jpg", "png", "pdf", "docx"]'),
  ('uploads.storage_path', '/var/www/uploads'),
  ('session.timeout', '1800')
ON CONFLICT ("key") DO UPDATE SET
  "value" = EXCLUDED."value",
  updated_at = now();
`,
}
