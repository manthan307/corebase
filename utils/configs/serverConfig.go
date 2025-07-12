package configs

import (
	"context"
	"encoding/json"
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
	CORSMaxAge       int
	TrustedProxies   []string
}

func InitConfig(port int, client *db.Client) Config {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Settings to pull from DB
	keys := []string{
		"server.read_timeout", "server.write_timeout", "server.idle_timeout",
		"cors.enabled", "cors.allow_methods", "cors.allow_headers", "cors.allow_credentials", "cors.max_age",
		"jwt.secret", "cors.allow_origins", "trusted_proxies", "allowed_hosts",
	}

	settings, err := client.Settings.GetBunch(ctx, keys)
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
	client.RedisClient.SetNX(ctx, "jwt:secret", jwtSecret, 0)

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
	client.RedisClient.SetNX(ctx, "trusted_proxies", trustedProxies, 0)

	// CORS Allow Origins
	corsOrigins := cfgMap["cors.allow_origins"]
	if corsOrigins == "" {
		host := helper.GetEnv("HOST", "")
		var origins []string

		if host != "" {
			origins = []string{host}
		} else {
			origins = []string{
				"http://localhost", "http://127.0.0.1", "http://[::1]",
				"http://localhost:3000", "http://127.0.0.1:3000",
			}
		}

		data, _ := json.Marshal(origins)
		corsOrigins = string(data)
		_ = client.Settings.Create(ctx, schema.SettingParams{
			Key:   "cors.allow_origins",
			Value: corsOrigins,
		})
	}
	client.RedisClient.SetNX(ctx, "cors:allow_origins", corsOrigins, 0)

	// Allowed Hosts (optional)
	allowedHosts := cfgMap["allowed_hosts"]
	if allowedHosts == "" {
		host := helper.GetEnv("HOST", "")
		if host != "" {
			data, _ := json.Marshal([]string{host})
			allowedHosts = string(data)
			_ = client.Settings.Create(ctx, schema.SettingParams{
				Key:   "allowed_hosts",
				Value: allowedHosts,
			})
			client.RedisClient.SetNX(ctx, "allowed_hosts", allowedHosts, 0)
		}
	} else {
		client.RedisClient.SetNX(ctx, "allowed_hosts", allowedHosts, 0)
	}

	// Final Config Object
	return Config{
		Port:             port,
		ReadTimeout:      helper.ParseDuration(cfgMap["server.read_timeout"], 15*time.Second),
		WriteTimeout:     helper.ParseDuration(cfgMap["server.write_timeout"], 15*time.Second),
		IdleTimeout:      helper.ParseDuration(cfgMap["server.idle_timeout"], 60*time.Second),
		CORSEnabled:      helper.ParseBool(cfgMap["cors.enabled"], true),
		CORSAllowMethods: helper.ParseStringSlice(cfgMap["cors.allow_methods"], []string{"GET", "POST", "PUT", "DELETE"}),
		CORSAllowHeaders: helper.ParseStringSlice(cfgMap["cors.allow_headers"], []string{"Content-Type", "Authorization"}),
		CORSAllowCreds:   helper.ParseBool(cfgMap["cors.allow_credentials"], true),
		CORSAllowOrigins: helper.ParseStringSlice(corsOrigins, helper.GetLocalAddresses()),
		CORSMaxAge:       helper.ParseInt(cfgMap["cors.max_age"], 720),
		TrustedProxies:   helper.ParseStringSlice(trustedProxies, helper.GetLocalAddresses()),
	}
}
