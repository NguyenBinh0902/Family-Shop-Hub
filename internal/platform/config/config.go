// Package config reads the application's settings from environment variables
// and validates them once at startup, so a misconfigured process fails fast
// with every problem listed instead of failing later on first use.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// LookupFunc has the signature of os.LookupEnv. Taking it as a parameter lets
// tests pass a map instead of mutating the real environment.
type LookupFunc func(key string) (string, bool)

// API holds the settings of the HTTP API process.
type API struct {
	Env             string // "development" or "production"
	HTTPAddr        string
	DatabaseURL     string // connects as the non-superuser application role
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

// Migration holds the settings of the migrate command.
type Migration struct {
	DatabaseURL string // connects as the owner role that may run DDL
}

// IsDevelopment reports whether the process runs in development mode.
func (c API) IsDevelopment() bool { return c.Env == "development" }

// LoadAPI reads the API settings. It collects every invalid or missing
// variable and returns them together with errors.Join.
func LoadAPI(lookup LookupFunc) (API, error) {
	r := reader{lookup: lookup}
	cfg := API{
		Env:             r.oneOf("APP_ENV", "production", "development", "production"),
		HTTPAddr:        r.string("HTTP_ADDR", ":8080"),
		DatabaseURL:     r.required("DATABASE_URL"),
		LogLevel:        r.logLevel("LOG_LEVEL", slog.LevelInfo),
		ShutdownTimeout: r.duration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
	return cfg, r.err()
}

// LoadMigration reads the settings of the migrate command.
func LoadMigration(lookup LookupFunc) (Migration, error) {
	r := reader{lookup: lookup}
	cfg := Migration{DatabaseURL: r.required("MIGRATION_DATABASE_URL")}
	return cfg, r.err()
}

// reader accumulates validation errors so one call reports all of them.
type reader struct {
	lookup LookupFunc
	errs   []error
}

func (r *reader) err() error { return errors.Join(r.errs...) }

func (r *reader) string(key, def string) string {
	if v, ok := r.lookup(key); ok && v != "" {
		return v
	}
	return def
}

func (r *reader) required(key string) string {
	v, ok := r.lookup(key)
	if !ok || v == "" {
		r.errs = append(r.errs, fmt.Errorf("%s is required", key))
	}
	return v
}

func (r *reader) oneOf(key, def string, allowed ...string) string {
	v := r.string(key, def)
	for _, a := range allowed {
		if v == a {
			return v
		}
	}
	r.errs = append(r.errs, fmt.Errorf("%s must be one of %v, got %q", key, allowed, v))
	return def
}

func (r *reader) duration(key string, def time.Duration) time.Duration {
	v, ok := r.lookup(key)
	if !ok || v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		r.errs = append(r.errs, fmt.Errorf("%s must be a positive duration such as 10s, got %q", key, v))
		return def
	}
	return d
}

func (r *reader) logLevel(key string, def slog.Level) slog.Level {
	v, ok := r.lookup(key)
	if !ok || v == "" {
		return def
	}
	var l slog.Level
	if err := l.UnmarshalText([]byte(v)); err != nil {
		r.errs = append(r.errs, fmt.Errorf("%s must be debug, info, warn or error, got %q", key, v))
		return def
	}
	return l
}
