package config

import (
	"os"
	"testing"
)

func TestConfigLoadDefaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()
	if cfg.AppPort != 8080 {
		t.Errorf("expected default AppPort 8080, got %d", cfg.AppPort)
	}
	if cfg.AppEnv != "development" {
		t.Errorf("expected default AppEnv development, got %s", cfg.AppEnv)
	}
	if !cfg.IsDevelopment() {
		t.Errorf("expected IsDevelopment() to be true")
	}
	if cfg.IsProduction() {
		t.Errorf("expected IsProduction() to be false")
	}
	if cfg.MinIOBuckets.Avatars != "avatars" {
		t.Errorf("expected default avatar bucket avatars, got %s", cfg.MinIOBuckets.Avatars)
	}
}

func TestConfigLoadWithCustomEnv(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("ALLOWED_ORIGINS", "https://temanmancing.id, https://app.temanmancing.id")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("ALLOWED_ORIGINS")
	}()

	cfg := Load()
	if cfg.AppPort != 9090 {
		t.Errorf("expected AppPort 9090, got %d", cfg.AppPort)
	}
	if !cfg.IsProduction() {
		t.Errorf("expected IsProduction() to be true")
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Fatalf("expected 2 allowed origins, got %d", len(cfg.AllowedOrigins))
	}
	if cfg.AllowedOrigins[0] != "https://temanmancing.id" {
		t.Errorf("expected origin https://temanmancing.id, got %s", cfg.AllowedOrigins[0])
	}
}
