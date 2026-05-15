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
