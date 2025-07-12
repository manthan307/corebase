package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/utils/helper"
)

func SecureHeaders(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		allowedHosts, err := client.RedisClient.Get(ctx,
			"allowed_hosts").Result()
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get allowed hosts"})
			return
		}

		var hosts []string
		if allowedHosts != "" {
			hosts = helper.ParseStringSlice(allowedHosts, []string{"localhost:8000"})
		}

		// Host header validation
		if len(hosts) > 0 {
			host := c.Request.Host
			valid := slices.Contains(hosts, host)
			if !valid {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
				return
			}
		}

		// Common secure headers
		c.Header("X-Frame-Options", "SAMEORIGIN") // allows iframe in dashboard if needed
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), fullscreen=(self)")

		// Don't cache sensitive data
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")

		// Content-Security-Policy (relaxed for dashboard, strict for API)
		if strings.HasPrefix(c.Request.URL.Path, "/_") {
			c.Header("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' 'unsafe-inline'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data:; "+
					"connect-src *; "+
					"font-src 'self';")
		} else {
			// Lock down API
			c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")
		}

		c.Next()
	}
}
