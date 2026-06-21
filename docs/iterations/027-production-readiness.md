# Iteration 027 - Production Readiness

## Goal

Prepare Banter Engine for long-running execution on Raspberry Pi before introducing AI components.

## Scope

### 1. Engine Runtime

Refactor orchestration from main.go into:

internal/engine/

Expose:

type Engine struct {}

func (e *Engine) Run(ctx context.Context) error

Responsibilities:
- fetch standings
- fetch matches
- fetch forecasts
- detect opportunities
- generate banter
- publish messages

### 2. Centralized Configuration

Create:

internal/config/

Expose:

type Config struct {
    MPPToken string
    DiscordWebhookURL string
    ChallengeID string
    SnapshotDir string
    PollInterval time.Duration
    DryRun bool
}

Configuration source:
- Environment variables

### 3. Graceful Shutdown

Support:
- CTRL+C
- SIGTERM
- systemd stop

Use signal.NotifyContext(...)

Requirements:
- stop scheduler
- flush logs
- exit cleanly

### 4. Dry Run Mode

Support:

banter-engine run
banter-engine dry-run

Dry run:
- no Discord call
- messages printed to stdout

### 5. Structured Logging

Create:

internal/logging/

Requirements:
- timestamp
- level
- component
- message

Use standard library only.

### 6. README Update

Document:
- installation
- environment variables
- dry run mode
- scheduler mode
- raspberry deployment

## Non Goals

Do NOT introduce:
- Ollama
- LLM
- embeddings
- vector databases

## Definition of Done

- go test ./... passes
- go vet ./... passes
- engine runs continuously
- graceful shutdown works
- dry run works
- configuration centralized
- logging centralized
- no regression
