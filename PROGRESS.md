# Progress

| Step | Status | Notes |
| --- | --- | --- |
| Update repo workflow docs | done | `AGENTS.md` now requires `PLAN.md` and `PROGRESS.md`. |
| Write active implementation plan | done | `PLAN.md` added for the marker pipeline v1. |
| Add local infrastructure config | done | Added `docker-compose.yml` plus env-backed config defaults. |
| Implement shared packages | done | Added domain, Postgres storage, Redis queue, and ffmpeg builder packages. |
| Implement API | done | Added health, marker create, marker get, and job get handlers. |
| Implement worker | done | Added Redis-driven processing with pending-job recovery. |
| Implement bot | done | Added transport-neutral handler, local stdin adapter, and HTTP marker client. |
| Implement CLI | done | Added HTTP-driven health, marker create/get, and job get commands. |
| Add tests and verify builds | done | `make tidy`, `make test`, and `make build` all succeed. |
| Expand local admin tooling | done | Added RedisInsight on `:5540` and pgAdmin on `:5050` with default local credentials. |
| Add shared port env file | done | Added root `.env` for API, Postgres, pgAdmin, Redis, and RedisInsight ports and endpoints. |
| Add Traefik stack config | done | Added external Traefik network wiring and default host/entrypoint/TLS env settings for pgAdmin and RedisInsight. |
| Fix admin/env defaults | done | Replaced reserved pgAdmin email domain and aligned app connection URLs with overridden Postgres/Redis host ports. |
| Add root `.env` config loading | done | `packages/config` now searches upward for the repo `.env` file and preserves explicit process env overrides. |
| Run end-to-end smoke test | done | Brought up Compose services, ran API on `:8088` due local `:8080` conflict, created a marker, and verified worker completion plus stored ffmpeg command. |
| Add basic README | done | Added setup, run, and smoke-test instructions for the current stack. |
| Refactor bot command handling | done | Added a reusable bot command framework with `!health` and `!marker` command handlers. |
| Expand bot health/help commands | done | Added `!help`, `!health:all`, `!health:bot`, `!health:server`, and typo-compatible `!heath:server`. |
| Add czg config | done | Added root `cz.config.cjs` with emoji commit types, Cuepoint scopes, and default prompt text. |
| Fix czg config shape | done | Switched `cz.config.cjs` to the documented `defineConfig` export so `czg` can load scopes correctly. |
| Containerize API and worker | done | Added shared Docker build config, Compose services, health-aware dependencies, and verified a containerized marker-to-worker flow via API on temporary host port `8090`. |
| Start Twitch integration | done | Added a first Twitch IRC adapter, env-driven bot transport selection, Twitch config fields, and tests for IRC parsing/reply behavior. |
| Fix Twitch bot config validation | done | Corrected `.env` channel wiring, added legacy `TWITCH_*` env fallbacks, and reject channel values that look like IRC addresses. |
| Add bot debug logging | done | Added leveled bot/Twitch diagnostics and changed `!help` to single-line chat-safe output. |
| Add bot uptime health output | done | Added uptime to bot health replies and made `!health` an alias for `!health:all`. |
