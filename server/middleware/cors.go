package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/utils/configs"
	"github.com/manthan307/corebase/utils/helper"
)

var (
	cacheMu        sync.Mutex
	cachedOrigins  []string
	lastFetchedAt  time.Time
	requestCounter int

	constMaxRequests     = 100             // fetch from Redis every 100 requests
	constRefreshInterval = 1 * time.Minute // or after 1 minute
)

func CORS(d configs.Config, client *db.Client) gin.HandlerFunc {
	if !d.CORSEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			now := time.Now()
			needRefresh := false

			cacheMu.Lock()

			requestCounter++
			if requestCounter >= constMaxRequests || now.Sub(lastFetchedAt) > constRefreshInterval {
				needRefresh = true
				requestCounter = 0
			}

			// Fetch from Redis if needed
			if needRefresh {
				go refreshOrigins(client, d)
			}

			origins := make([]string, len(cachedOrigins))
			copy(origins, cachedOrigins)
			cacheMu.Unlock()

			// Match origin
			for _, o := range origins {
				if o == origin {
					return true
				}
			}
			return false
		},
		AllowMethods:     d.CORSAllowMethods,
		AllowHeaders:     d.CORSAllowHeaders,
		AllowCredentials: d.CORSAllowCreds,
		MaxAge:           time.Duration(d.CORSMaxAge) * time.Minute,
	})
}

func refreshOrigins(client *db.Client, d configs.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	raw, err := client.RedisClient.Get(ctx, "cors:allow_origins").Result()

	var origins []string
	if err != nil || raw == "" {
		origins = d.CORSAllowOrigins
	} else if err := helper.ParseJSON(raw, &origins); err != nil || len(origins) == 0 {
		origins = d.CORSAllowOrigins
	}

	cacheMu.Lock()
	cachedOrigins = origins
	lastFetchedAt = time.Now()
	cacheMu.Unlock()
}
