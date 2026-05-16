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
| Replace plan with control-center v1 | done | `PLAN.md` now describes the hybrid Next.js + Go control-center implementation and follow-up work. |
| Scaffold `apps/web` workspace | done | Added `pnpm-workspace.yaml`, root web scripts, and a new Next.js/Tailwind/shadcn-style app scaffold under `apps/web`. |
| Add control-center database schema | done | `packages/database` now applies embedded SQL migrations for operator profiles, bot settings, plans, segments, sessions, timeline markers, and automation jobs. |
| Add shared runtime domain types | done | `packages/stream` now defines control-center plans, segments, sessions, title formatting, markers, and automation job types. |
| Add API runtime endpoints | done | `apps/api` now exposes internal session sync, title control, segment control, and timeline marker endpoints guarded by the shared internal token header. |
| Extend bot runtime commands | done | `apps/bot` now supports title commands, viewing/react commands, next-segment flow, timeline markers, and periodic runtime sync. |
| Add worker YouTube metadata jobs | done | `apps/worker` now claims pending `youtube_metadata_sync` jobs and fills segment metadata from YouTube oEmbed. |
| Add web auth, settings, plans, and overlays | done | `apps/web` now includes login, dashboard, plan pages, Twitch settings, bot settings, internal Twitch proxy routes, and public overlay routes/pages. |
| Verify Go control-center changes | done | Ran `gofmt`, package/app `go test`, and app `go build` with writable `/tmp` Go caches. |
| Verify frontend dependency install/build | done | `pnpm --dir apps/web build` now succeeds after marking auth-backed pages dynamic. |
| Fix login page build collection | done | Moved `/login` to a client-rendered shell so Next no longer executes Better Auth and DB bootstrap during page data collection. |
| Fix active react segment resolution | done | `ResolveActiveSegment` now prefers `session.current_segment_id`, so manual/live segment state stays visible to `!react` and `!watching`. |
| Commit Better Auth generated schema migration | done | Added repo migration `002_better_auth.sql` for Better Auth core tables and web-side schema bootstrap for first-run auth requests. |
| Refactor web shell into sidebar dashboard | done | Reworked `apps/web` around a persistent sidebar shell with mobile nav fallbacks and a dashboard-style layout. |
| Add toggleable dark mode | done | Added a client-side theme toggle with persisted local storage state and dark-mode-aware UI primitives/styles. |
| Convert small-option text inputs to selectors | done | Plan status, segment type, and timing mode now use select controls instead of free-text inputs. |
| Add field help text | done | Added explanatory subheadings to planning and settings forms, including elapsed offset, duration, and title format behavior. |
| Add dev forced-live settings page | done | Added `/settings/dev` plus `dev_user_settings` storage and stream-state overrides for local runtime testing. |
| Add Twitch Helix handler structure | done | Added markers, videos, and schedule handler modules plus internal route bridges near the existing Twitch channel/stream bridge code. |
