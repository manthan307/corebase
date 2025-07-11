package configs

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/manthan307/corebase/db"
	"github.com/manthan307/corebase/db/schema"
	"github.com/manthan307/corebase/utils/helper"
)

type Config struct {
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	CORSEnabled      bool
	CORSAllowMethods []string
	CORSAllowHeaders []string
	CORSAllowCreds   bool
	CORSAllowOrigins []string
	JWTSecret        string
	RateLimit        int
	TrustedProxies   []string
}

func InitConfig(port int, client *db.Client) Config {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings, err := client.Settings.GetBunch(ctx, []string{
		"server.read_timeout",
		"server.write_timeout",
		"server.idle_timeout",
		"cors.enabled",
		"cors.allow_methods",
		"cors.allow_headers",
		"cors.allow_credentials",
		"rate_limit",
		"jwt.secret",
		"cors.allow_origins",
		"trusted_proxies",
	})
	if err != nil {
		panic(err)
	}

	cfgMap := make(map[string]string)
	for _, s := range settings {
		cfgMap[s.Key] = s.Value
	}

	// JWT Secret
	jwtSecret := cfgMap["jwt.secret"]
	if jwtSecret == "" {
		jwtSecret = helper.GenerateRandomString(32)
		_ = client.Settings.Create(ctx, schema.SettingParams{
			Key:   "jwt.secret",
			Value: jwtSecret,
		})
	}

	// Trusted Proxies
	trustedProxies := cfgMap["trusted_proxies"]
	if trustedProxies == "" {
		proxies := helper.GetLocalAddresses()
		data, _ := json.Marshal(proxies)
		trustedProxies = string(data)
		_ = client.Settings.Create(ctx, schema.SettingParams{
			Key:   "trusted_proxies",
			Value: trustedProxies,
		})
	}

	// CORS Allow Origins
	corsOrigins := cfgMap["cors.allow_origins"]
	if corsOrigins == "" {
		host := os.Getenv("HOST")
		if host != "" {
			origins := []string{host}
			data, _ := json.Marshal(origins)
			corsOrigins = string(data)
			_ = client.Settings.Create(ctx, schema.SettingParams{
				Key:   "cors.allow_origins",
				Value: corsOrigins,
			})
		}
	}

	return Config{
		Port:             port,
		ReadTimeout:      helper.ParseDuration(cfgMap["server.read_timeout"], 15*time.Second),
		WriteTimeout:     helper.ParseDuration(cfgMap["server.write_timeout"], 15*time.Second),
		IdleTimeout:      helper.ParseDuration(cfgMap["server.idle_timeout"], 60*time.Second),
		CORSEnabled:      helper.ParseBool(cfgMap["cors.enabled"], true),
		CORSAllowMethods: helper.ParseStringSlice(cfgMap["cors.allow_methods"], []string{"GET", "POST", "PUT", "DELETE"}),
		CORSAllowHeaders: helper.ParseStringSlice(cfgMap["cors.allow_headers"], []string{"Content-Type", "Authorization"}),
		CORSAllowCreds:   helper.ParseBool(cfgMap["cors.allow_credentials"], true),
		RateLimit:        helper.ParseInt(cfgMap["rate_limit"], 10),
		JWTSecret:        jwtSecret,
		CORSAllowOrigins: helper.ParseStringSlice(corsOrigins, helper.GetLocalAddresses()),
		TrustedProxies:   helper.ParseStringSlice(trustedProxies, helper.GetLocalAddresses()),
	}
}
