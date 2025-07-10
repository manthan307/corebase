package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/utils/configs"
	"github.com/manthan307/corebase/utils/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("server",
	fx.Provide(ProvideServer),
	fx.Invoke(StartServer),
)

func ProvideServer(d configs.Config, log *zap.Logger) *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprint(d.Port)
	}

	router := gin.New()

	router.Use(logger.LoggerMiddleware(log))
	router.Use(gin.Recovery())

	return &http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: router.Handler(),
	}
}

func StartServer(lc fx.Lifecycle, d configs.Config, router *http.Server, log *zap.Logger) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprint(d.Port)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("🚀 Listing on port " + port)
				err := router.ListenAndServe()
				if err != nil {
					log.Error("Error starting server:", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping server.")
			Shutdownctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			router.Shutdown(Shutdownctx)
			return nil
		},
	},
	)
}
