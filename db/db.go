package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manthan307/corebase/db/client"
	"github.com/manthan307/corebase/utils/helper"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"db",
	fx.Provide(InitPostgresDataBase, NewClient),
	fx.Invoke(RunMigration),
)

func InitPostgresDataBase(lc fx.Lifecycle, log *zap.Logger) *pgxpool.Pool {
	dbURL := helper.GetEnv("PG_URL", "")
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
	Db          *pgxpool.Pool
	Settings    client.SettingRepo
	Admin       client.AdminsRepo
	RedisClient *redis.Client
}

func NewClient(pool *pgxpool.Pool, log *zap.Logger) *Client {
	return &Client{
		Db:          pool,
		Settings:    client.NewSettingRepo(pool),
		RedisClient: NewRedisClient(log),
		Admin:       client.NewAdminRepo(pool),
	}
}
