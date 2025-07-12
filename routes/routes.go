package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/routes/v1/admin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"routes",
	fx.Provide(ProvideV1Group),
	fx.Invoke(
		admin.AdminRoutes,
	),
)

func ProvideV1Group(server *http.Server, log *zap.Logger) *gin.RouterGroup {
	// cast back the router
	engine, ok := server.Handler.(*gin.Engine)
	if ok {
		routerGroup := engine.Group("/api/v1")
		return routerGroup
	}

	log.Fatal("msg", zap.Error(errors.New("cannot cast to gin.Engine")))
	return nil
}
