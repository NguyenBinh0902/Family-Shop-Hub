package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func mapLookup(m map[string]string) LookupFunc {
	return func(k string) (string, bool) {
		v, ok := m[k]
		return v, ok
	}
}

func TestLoadAPI_Defaults(t *testing.T) {
	cfg, err := LoadAPI(mapLookup(map[string]string{"DATABASE_URL": "postgres://x"}))
	if err != nil {
		t.Fatalf("LoadAPI: %v", err)
	}
	want := API{
		Env:             "production",
		HTTPAddr:        ":8080",
		DatabaseURL:     "postgres://x",
		LogLevel:        slog.LevelInfo,
		ShutdownTimeout: 10 * time.Second,
	}
	if cfg != want {
		t.Errorf("cfg = %+v, want %+v", cfg, want)
	}
}

func TestLoadAPI_ReadsValues(t *testing.T) {
	cfg, err := LoadAPI(mapLookup(map[string]string{
		"APP_ENV":          "development",
		"HTTP_ADDR":        ":9000",
		"DATABASE_URL":     "postgres://x",
		"LOG_LEVEL":        "debug",
		"SHUTDOWN_TIMEOUT": "3s",
	}))
	if err != nil {
		t.Fatalf("LoadAPI: %v", err)
	}
	if !cfg.IsDevelopment() || cfg.HTTPAddr != ":9000" || cfg.LogLevel != slog.LevelDebug || cfg.ShutdownTimeout != 3*time.Second {
		t.Errorf("unexpected cfg %+v", cfg)
	}
}

func TestLoadAPI_ReportsAllErrorsAtOnce(t *testing.T) {
	_, err := LoadAPI(mapLookup(map[string]string{
		"APP_ENV":          "staging",
		"LOG_LEVEL":        "loud",
		"SHUTDOWN_TIMEOUT": "soon",
	}))
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, key := range []string{"DATABASE_URL", "APP_ENV", "LOG_LEVEL", "SHUTDOWN_TIMEOUT"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error %q does not mention %s", err, key)
		}
	}
}

func TestLoadMigration_RequiresURL(t *testing.T) {
	if _, err := LoadMigration(mapLookup(nil)); err == nil {
		t.Error("expected an error when MIGRATION_DATABASE_URL is missing")
	}
}
