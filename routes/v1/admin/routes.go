package admin

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"go.uber.org/zap"
)

func AdminRoutes(rg *gin.RouterGroup, log *zap.Logger, client *db.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	TotalCount, err := client.Admin.TotalCount(ctx)
	if err != nil {
		log.Error("Error getting total count", zap.Error(err))
	}

	if TotalCount <= 0 {
		url, err := OneTimeURL(ctx, client, "superadmin", 0)
		if err != nil {
			log.Error("Error generating one time url", zap.Error(err))
		}
		log.Warn("One time url", zap.String("url", url))
	}

	admin := rg.Group("/admin")

	//Routes
	admin.POST("/create", CreateAdmin(log, client))
}
