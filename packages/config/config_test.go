package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLoadReadsRootDotEnv(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
		loadEnvOnce = sync.Once{}
	})

	root := t.TempDir()
	appDir := filepath.Join(root, "apps", "api")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("API_ADDR=:9999\nREDIS_DB=7\nBOT_TWITCH_TLS=false\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Chdir(appDir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	_ = os.Unsetenv("API_ADDR")
	_ = os.Unsetenv("REDIS_DB")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.API.Addr != ":9999" {
		t.Fatalf("API.Addr = %q", cfg.API.Addr)
	}
	if cfg.Redis.DB != 7 {
		t.Fatalf("Redis.DB = %d", cfg.Redis.DB)
	}
	if cfg.Bot.Twitch.UseTLS {
		t.Fatal("expected Bot.Twitch.UseTLS to be false")
	}
}

func TestLoadDoesNotOverrideExistingEnv(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
		loadEnvOnce = sync.Once{}
	})

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("API_ADDR=:9999\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	t.Setenv("API_ADDR", ":1234")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.API.Addr != ":1234" {
		t.Fatalf("API.Addr = %q", cfg.API.Addr)
	}
}

func TestLoadSupportsLegacyTwitchEnvKeys(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
		loadEnvOnce = sync.Once{}
	})

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte("TWITCH_BOT_USERNAME=botuser\nTWITCH_OAUTH_TOKEN=oauth:token\nTWITCH_CHANNEL=channelname\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	_ = os.Unsetenv("BOT_TWITCH_USERNAME")
	_ = os.Unsetenv("BOT_TWITCH_OAUTH_TOKEN")
	_ = os.Unsetenv("BOT_TWITCH_CHANNEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Bot.Twitch.Username != "botuser" {
		t.Fatalf("Username = %q", cfg.Bot.Twitch.Username)
	}
	if cfg.Bot.Twitch.OAuthToken != "oauth:token" {
		t.Fatalf("OAuthToken = %q", cfg.Bot.Twitch.OAuthToken)
	}
	if cfg.Bot.Twitch.Channel != "channelname" {
		t.Fatalf("Channel = %q", cfg.Bot.Twitch.Channel)
	}
}
