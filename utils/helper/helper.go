package helper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a string into time.Duration. Defaults to fallback on error.
func ParseDuration(s string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}

// ParseBool parses a string into bool. Defaults to fallback on error.
func ParseBool(s string, fallback bool) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return fallback
	}
	return b
}

// ParseInt parses a string into int. Defaults to fallback on error.
func ParseInt(s string, fallback int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return i
}

// ParseStringSlice parses a JSON string slice (e.g. '["GET","POST"]')
// Returns fallback if unmarshal fails.
func ParseStringSlice(s string, fallback []string) []string {
	var out []string
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return fallback
	}
	return out
}

// GetEnv returns the value of an environment variable or a fallback
func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// GetEnvSlice splits an environment variable by comma or returns a fallback
func GetEnvSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return fallback
}

// GenerateRandomString returns a cryptographically secure random base64 string
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to non-crypto random
		for i := range b {
			val, _ := rand.Int(rand.Reader, big.NewInt(256))
			b[i] = byte(val.Int64())
		}
	}
	return base64.URLEncoding.EncodeToString(b)
}

// GetLocalAddresses returns local IPs like 127.0.0.1 and ::1
func GetLocalAddresses() []string {
	return []string{"127.0.0.1", "::1", "localhost"}
}

func ParseJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
