package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"Effective-Mobile/internal/config"
)

func chdirProjectRoot(t *testing.T) {
	t.Helper()
	wd, _ := os.Getwd()
	root := filepath.Join(wd, "../..")
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
}

func TestLoad_Default(t *testing.T) {
	chdirProjectRoot(t)
	t.Setenv("APP_ENV", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.Server.Port == 0 {
		t.Error("expected non‑zero port from config file")
	}
	if cfg.Enrichment.Timeout == 0*time.Second {
		t.Error("timeout not parsed")
	}
}

func TestLoad_UnknownEnv(t *testing.T) {
	chdirProjectRoot(t)
	t.Setenv("APP_ENV", "unknown-env-file")
	if _, err := config.Load(); err == nil {
		t.Error("expected error for missing env‑specific file")
	}
}
