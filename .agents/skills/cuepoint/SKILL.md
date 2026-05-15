---
name: cuepoint
description: Use when working in the Cuepoint Go monorepo.
---

# Cuepoint workflow

Before editing:
- inspect `go.work`
- inspect the relevant module's `go.mod`
- check existing package boundaries

Common commands:

```bash
go work sync
make tidy
make test
make build
```


When adding a feature:

shared domain types go in packages/stream
ffmpeg command generation goes in packages/ffmpeg
chatbot command parsing goes in apps/bot
API handlers go in apps/api
async processing goes in apps/worker

Do not hardcode Twitch, YouTube, or ffmpeg credentials.
