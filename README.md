# Cuepoint

Cuepoint is a Go monorepo for stream scheduling, chatbot markers, VOD download links, and ffmpeg segment commands.

## Services

- `apps/api` exposes the HTTP API.
- `apps/worker` processes queued jobs and stores generated ffmpeg commands.
- `apps/bot` provides a transport-agnostic bot framework with local stdin and Twitch IRC adapters.
- `apps/cli` provides local admin and smoke-test commands.
- `packages/*` contains shared config, domain, storage, queue, and ffmpeg logic.

## Local requirements

- Go `1.26.2`
- Docker with Compose

## Local stack

Start the full local stack:

```bash
docker compose up -d
```

Default local ports come from [.env](/home/hamish/Development/MarsuvesVex/public-projects/cuepoint/.env):

- API host port: `8088`
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

This now includes:

- `api`
- `worker`
- `postgres`
- `pgadmin`
- `redis`
- `redisinsight`

The `api` and `worker` containers override the host-style `.env` values with container-native service addresses for Postgres and Redis.
The API container listens internally on `API_CONTAINER_PORT` and is published on `API_PORT`, so host port conflicts can be resolved without changing the container wiring.

## Run the apps manually

Each binary still loads values from the root `.env` file automatically, so you can run them outside Compose when needed.

Start the API manually:

```bash
./apps/api/api
```

Start the worker manually:

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

To use Twitch chat instead of stdin, set these values in `.env` before running the bot:

```bash
BOT_TRANSPORT=twitch
BOT_LOG_LEVEL=debug
BOT_TWITCH_USERNAME=your_bot_username
BOT_TWITCH_OAUTH_TOKEN=oauth:your_token
BOT_TWITCH_CHANNEL=your_channel
```

Bot log levels:

- `debug` shows IRC connection, parsed commands, ignored IRC lines, and sent replies
- `info` shows startup, transport selection, connect/join, and shutdown events
- `warn` shows command handling problems
- `error` shows only failures

Example bot commands:

```text
!help
!health
!health:all
!health:bot
!health:server
!marker vod-source.mp4 intro 00:00:10
```

## Current workflow

1. Start Compose services.
2. Use the CLI or bot against the running API.
3. Create a marker through the CLI or bot.
4. Fetch the job and confirm the stored ffmpeg command.

## Admin tools

- pgAdmin: `http://localhost:5054`
  - email: `admin@cuepoint.app`
  - password: `cuepoint`
- RedisInsight: `http://localhost:5541`
