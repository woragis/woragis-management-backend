package middleware

import (
	"os"
	"strings"
)

type Config struct {
	CORSOrigins          []string
	CORSAllowCredentials bool
	MetricsEnabled       bool
}

func normalizeOrigin(raw string) string {
	o := strings.TrimSpace(raw)
	o = strings.Trim(o, `"'`)
	return o
}

func parseCORSOrigins(value string) []string {
	var origins []string
	for _, part := range strings.Split(value, ",") {
		if o := normalizeOrigin(part); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}

func LoadConfigFromEnv() Config {
	origins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
	if v := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); v != "" {
		origins = parseCORSOrigins(v)
	}
	allowCreds := true
	if v := strings.TrimSpace(os.Getenv("CORS_ALLOW_CREDENTIALS")); v == "0" || strings.EqualFold(v, "false") {
		allowCreds = false
	}
	metricsOn := true
	if v := strings.TrimSpace(os.Getenv("METRICS_ENABLED")); v == "0" || strings.EqualFold(v, "false") {
		metricsOn = false
	}
	return Config{
		CORSOrigins:          origins,
		CORSAllowCredentials: allowCreds,
		MetricsEnabled:       metricsOn,
	}
}
