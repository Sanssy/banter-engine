package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultChallengeID  = "mpp_challenge_UDKDDH27"
	defaultSnapshotDir  = "data"
	defaultPollInterval = 5 * time.Minute
)

type Config struct {
	MPPToken          string
	DiscordWebhookURL string
	ChallengeID       string
	SnapshotDir       string
	PollInterval      time.Duration
	DryRun            bool
}

func Load() (Config, error) {
	cfg := Config{
		MPPToken:          os.Getenv("MPP_TOKEN"),
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		ChallengeID:       envOrDefault("CHALLENGE_ID", defaultChallengeID),
		SnapshotDir:       envOrDefault("SNAPSHOT_DIR", defaultSnapshotDir),
		PollInterval:      defaultPollInterval,
	}

	if value := os.Getenv("POLL_INTERVAL"); value != "" {
		interval, err := time.ParseDuration(value)
		if err != nil || interval <= 0 {
			return Config{}, fmt.Errorf("POLL_INTERVAL must be a positive duration")
		}
		cfg.PollInterval = interval
	}
	if value := os.Getenv("DRY_RUN"); value != "" {
		dryRun, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("DRY_RUN must be a boolean: %w", err)
		}
		cfg.DryRun = dryRun
	}

	if cfg.MPPToken == "" {
		return Config{}, fmt.Errorf("MPP_TOKEN environment variable is required")
	}
	if !cfg.DryRun && cfg.DiscordWebhookURL == "" {
		return Config{}, fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
	}
	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
