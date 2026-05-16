# Cuepoint Control Center v1

## Summary

Implement a hybrid control center:

- `apps/web` owns the operator UI, Better Auth, Twitch account linking, server actions, and OBS overlay routes.
- `apps/api` owns bot/runtime orchestration for session sync, segment advancement, title automation, and marker writes.
- `apps/bot` owns chat command parsing plus periodic live-session sync.
- `apps/worker` owns asynchronous YouTube metadata enrichment.
- Postgres remains the shared state store for plans, segments, runtime sessions, markers, settings, and automation jobs.

## Implemented In This Pass

- Added a new `apps/web` Next.js workspace scaffold with App Router, Tailwind, and shadcn-style UI primitives.
- Added Better Auth wiring, Twitch sign-in entrypoint, Twitch profile sync, internal Twitch proxy routes, and public overlay data routes.
- Added control-center data tables through embedded SQL migrations in `packages/database`.
- Added shared stream/runtime domain types for plans, segments, sessions, markers, settings, and automation jobs.
- Added API runtime endpoints for session sync, title apply/restore/toggle, title format management, segment start/advance, and timeline markers.
- Added bot commands:
  - `!settitle`
  - `!restoretitle`
  - `!toggletitles`
  - `!titleformat`
  - `!viewtitleformat`
  - `!resettitleformat`
  - `!watching`
  - `!react`
  - `!nextsegment`
  - `!marker <label> [end]` for timeline markers alongside the legacy ffmpeg marker flow
- Added worker support for `youtube_metadata_sync` automation jobs via YouTube oEmbed.

## Remaining Next Steps

- Install web dependencies and run frontend verification (`pnpm install`, `pnpm --dir apps/web build`, and lint/type-check passes).
- Add committed Better Auth table migrations once the exact generated schema is locked from a dependency install.
- Expand server-side validation in the web app for plan/segment forms and replace free-text status/type inputs with constrained UI controls.
- Add integration tests for internal API runtime routes, title formatting behavior, marker pairing, and worker YouTube metadata processing.
- Add overlay styling refinements and a stronger public overlay token story if the current `overlay_public_id` needs rotation or revocation flows.
