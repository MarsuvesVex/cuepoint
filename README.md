# Cuepoint

Cuepoint is a Go monorepo for stream scheduling, chatbot markers, VOD download links, and ffmpeg segment commands.

## Services

- `apps/api` exposes the HTTP API.
- `apps/worker` processes queued jobs and stores generated ffmpeg commands.
- `apps/bot` provides a transport-agnostic bot framework with a local stdin adapter today.
- `apps/cli` provides local admin and smoke-test commands.
- `packages/*` contains shared config, domain, storage, queue, and ffmpeg logic.

## Local requirements

- Go `1.26.2`
- Docker with Compose

## Local stack

Start the local infrastructure:

```bash
docker compose up -d
```

Default local ports come from [.env](/home/hamish/Development/MarsuvesVex/public-projects/cuepoint/.env):

- API: `8080`
- Postgres: `5439`
- pgAdmin: `5054`
- Redis: `6378`
- RedisInsight: `5541`

Traefik-related hostnames and network settings also live in `.env`.

## Build and test

From the repo root:

```bash
make tidy
make test
make build
```

## Run the apps

Each binary loads values from the root `.env` file automatically.

Start the API:

```bash
./apps/api/api
```

Start the worker:

```bash
./apps/worker/worker
```

Use the CLI:

```bash
./apps/cli/cuepoint health
./apps/cli/cuepoint marker create --stream vod-source.mp4 --label intro --timestamp 00:00:10
```

Run the local bot adapter:

```bash
./apps/bot/bot
```

Example bot commands:

```text
!help
!health:all
!health:bot
!health:server
!marker vod-source.mp4 intro 00:00:10
```

## Current workflow

1. Start Compose services.
2. Run the API and worker.
3. Create a marker through the CLI or bot.
4. Fetch the job and confirm the stored ffmpeg command.

## Admin tools

- pgAdmin: `http://localhost:5054`
  - email: `admin@cuepoint.app`
  - password: `cuepoint`
- RedisInsight: `http://localhost:5541`
