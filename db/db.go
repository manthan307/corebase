package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manthan307/corebase/db/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"db",
	fx.Provide(InitPostgresDataBase, NewClient),
	fx.Invoke(RunMigration),
)

func InitPostgresDataBase(lc fx.Lifecycle, log *zap.Logger) *pgxpool.Pool {
	dbURL := os.Getenv("PG_URL")
	if dbURL == "" {
		log.Fatal("PG_URL is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database", zap.Error(err))
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Unable to connect to database", zap.Error(err))
	}

	log.Info("✅ Connected to database")

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("🔌 Disconnecting from database")
				db.Close()
				return nil
			},
		},
	)
	return db
}

type Client struct {
	Db       *pgxpool.Pool
	Settings client.SettingRepo
}

func NewClient(pool *pgxpool.Pool) *Client {
	return &Client{
		Db:       pool,
		Settings: client.NewSettingRepo(pool),
	}
}
