# Marker Pipeline v1 Plan

## Summary

Implement the first end-to-end vertical slice for Cuepoint using Postgres and Redis from day one. The initial workflow creates stream markers through the API, persists them in Postgres, enqueues background jobs in Redis, and lets the worker turn those jobs into stored ffmpeg command output. The bot remains transport-agnostic so a Twitch adapter can be added next without changing command logic.

## Key Changes

- Add a root `docker-compose.yml` with Postgres and Redis for local development.
- Add shared packages for config, database access, Redis queueing, stream domain services, and ffmpeg command generation.
- Implement four binaries:
  - API with health, marker creation, marker lookup, and job lookup
  - Worker with Redis-driven job processing and startup requeue of pending jobs
  - Bot with adapter interfaces plus a local stdin adapter
  - CLI for health and marker/job inspection flows
- Keep app-specific wiring under `apps/*/internal`.

## Public Interfaces

- HTTP API:
  - `GET /healthz`
  - `POST /markers`
  - `GET /markers/{id}`
  - `GET /jobs/{id}`
- Bot command:
  - `!marker <stream> <label> <timestamp>`
- CLI commands:
  - `cuepoint health`
  - `cuepoint marker create --stream ... --label ... --timestamp ...`
  - `cuepoint marker get --id ...`
  - `cuepoint job get --id ...`

## Test Plan

- Unit tests for ffmpeg command generation and stream service behavior.
- Handler tests for API validation and successful marker creation.
- Bot parser tests for transport-neutral command handling.
- Build verification for all four binaries and package tests under `make test`.
