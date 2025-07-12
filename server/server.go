package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/routes"
	"github.com/manthan307/corebase/server/middleware"
	"github.com/manthan307/corebase/utils/configs"
	"github.com/manthan307/corebase/utils/helper"
	"github.com/manthan307/corebase/utils/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("server",
	fx.Provide(ProvideServer),
	fx.Invoke(StartServer),
	routes.Module,
)

func ProvideServer(d configs.Config, log *zap.Logger, client *db.Client) *http.Server {
	port := d.Port
	if envPort := helper.GetEnv("PORT", "8000"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(logger.LoggerMiddleware(log))
	router.Use(gin.Recovery())

	router.SetTrustedProxies(d.TrustedProxies)

	router.Use(middleware.SecureHeaders(client))
	router.Use(middleware.CORS(d, client))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  d.ReadTimeout,
		WriteTimeout: d.WriteTimeout,
		IdleTimeout:  d.IdleTimeout,
	}
}

func StartServer(lc fx.Lifecycle, d configs.Config, server *http.Server, log *zap.Logger, shutdown fx.Shutdowner) {
	// cast back the router
	engine, ok := server.Handler.(*gin.Engine)
	if ok {
		//TODO
		engine.POST("/internal/shutdown", func(c *gin.Context) {
			secret := helper.GetEnv("SHUTDOWN_SECRET", "")
			if secret != "" && c.GetHeader("X-Secret") != secret {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			// go func() {
			// 	log.Warn("🧨 Shutdown triggered from /internal/shutdown")
			// 	if err := shutdown.Shutdown(); err != nil {
			// 		log.Error("Failed to trigger shutdown", zap.Error(err))
			// 	}
			// }()

			c.JSON(http.StatusOK, gin.H{"status": "shutting down"})
		})
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info(fmt.Sprintf("🚀 Listening on port %s", server.Addr))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("❌ Server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("🛑 Shutting down server...")
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return server.Shutdown(shutdownCtx)
		},
	})
}
