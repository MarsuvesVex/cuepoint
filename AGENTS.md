# AGENTS.md

## Project

Cuepoint is a Go monorepo for stream scheduling, chatbot markers, VOD download links, and ffmpeg segment commands.

## Layout

- `apps/api` — HTTP API
- `apps/bot` — stream chatbot
- `apps/worker` — VOD/ffmpeg background jobs
- `apps/cli` — local admin CLI
- `packages/*` — shared Go packages

## Commands

Use these from the repository root:

```bash
make tidy
make test
make build
```

For one module:

```bash
go test ./...
go test -C apps/api ./...
go test -C packages/stream ./...
```

## Rules

Keep app-specific code inside apps/*/internal.
Put shared types and logic in packages/*.
Do not put secrets in examples.
Prefer small, testable Go packages.
When changing shared packages, check affected apps.
Use context.Context for IO, DB, HTTP, and worker operations.
Write all approved implementation plans into `PLAN.md`.
Update `PROGRESS.md` as each implementation step is completed, using explicit statuses.
