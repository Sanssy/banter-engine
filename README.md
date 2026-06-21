# Banter Engine

Banter Engine polls Mon Petit Prono data, detects noteworthy football prediction events, generates deterministic banter, and publishes it to Discord.

## Installation

Go 1.26 or later is required.

```sh
go build -o banter-engine ./cmd/banter-engine
```

Run the binary from the repository root so it can read `resources/opportunities.json`.

## Configuration

Configuration is read from environment variables:

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `MPP_TOKEN` | yes | none | Mon Petit Prono API token |
| `DISCORD_WEBHOOK_URL` | except in dry-run mode | none | Discord webhook receiving messages |
| `CHALLENGE_ID` | no | `mpp_challenge_UDKDDH27` | MPP challenge to monitor |
| `SNAPSHOT_DIR` | no | `data` | Runtime snapshot directory |
| `POLL_INTERVAL` | no | `5m` | Positive Go duration between polls |
| `DRY_RUN` | no | `false` | Disable Discord and print messages to stdout |

## Running

Continuous scheduler mode:

```sh
MPP_TOKEN=... DISCORD_WEBHOOK_URL=... ./banter-engine run
```

Dry-run mode performs the same polling and detection without contacting Discord:

```sh
MPP_TOKEN=... ./banter-engine dry-run
```

The process handles `SIGINT` and `SIGTERM`, stops its scheduler, and exits cleanly.

## Raspberry Pi

Build directly on the Raspberry Pi, or cross-compile for a 64-bit Raspberry Pi:

```sh
GOOS=linux GOARCH=arm64 go build -o banter-engine ./cmd/banter-engine
```

Place the binary, the `resources` directory, and a writable snapshot directory under `/opt/banter-engine`. A minimal systemd unit is:

```ini
[Unit]
Description=Banter Engine
After=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/banter-engine
ExecStart=/opt/banter-engine/banter-engine run
EnvironmentFile=/etc/banter-engine.env
Restart=on-failure
User=banter-engine

[Install]
WantedBy=multi-user.target
```

After creating `/etc/banter-engine.env`, enable the service with `sudo systemctl enable --now banter-engine`. `systemctl stop banter-engine` uses the graceful `SIGTERM` shutdown path.
