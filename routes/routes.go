package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"routes",
	fx.Invoke(Routes),
)

func Routes(server *http.Server, log *zap.Logger, client *db.Client) {
	// cast back the router
	engine, ok := server.Handler.(*gin.Engine)
	if ok {
		engine.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

}
