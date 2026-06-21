package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("MPP_TOKEN", "token")
	t.Setenv("DISCORD_WEBHOOK_URL", "https://discord.example/webhook")
	t.Setenv("CHALLENGE_ID", "challenge")
	t.Setenv("SNAPSHOT_DIR", "state")
	t.Setenv("POLL_INTERVAL", "30s")
	t.Setenv("DRY_RUN", "false")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.MPPToken != "token" || cfg.DiscordWebhookURL == "" || cfg.ChallengeID != "challenge" {
		t.Fatalf("Load() config = %+v", cfg)
	}
	if cfg.SnapshotDir != "state" || cfg.PollInterval != 30*time.Second || cfg.DryRun {
		t.Fatalf("Load() config = %+v", cfg)
	}
}

func TestLoadDryRunDoesNotRequireDiscord(t *testing.T) {
	t.Setenv("MPP_TOKEN", "token")
	t.Setenv("DISCORD_WEBHOOK_URL", "")
	t.Setenv("DRY_RUN", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.DryRun {
		t.Fatal("Load() DryRun = false, want true")
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	t.Setenv("MPP_TOKEN", "")
	t.Setenv("DRY_RUN", "true")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "MPP_TOKEN") {
		t.Fatalf("Load() error = %v, want MPP_TOKEN error", err)
	}

	t.Setenv("MPP_TOKEN", "token")
	t.Setenv("POLL_INTERVAL", "never")
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "POLL_INTERVAL") {
		t.Fatalf("Load() error = %v, want POLL_INTERVAL error", err)
	}
}
