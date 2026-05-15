.PHONY: tidy build test

GOENV = GOCACHE=/tmp/cuepoint-go-build GOMODCACHE=/tmp/cuepoint-go-mod GOFLAGS=-buildvcs=false

tidy:
	$(GOENV) go work sync
	$(GOENV) go mod tidy -C apps/api
	$(GOENV) go mod tidy -C apps/bot
	$(GOENV) go mod tidy -C apps/worker
	$(GOENV) go mod tidy -C apps/cli
	$(GOENV) go mod tidy -C packages/config
	$(GOENV) go mod tidy -C packages/database
	$(GOENV) go mod tidy -C packages/events
	$(GOENV) go mod tidy -C packages/ffmpeg
	$(GOENV) go mod tidy -C packages/stream

build:
	$(GOENV) go build -C apps/api ./cmd/api
	$(GOENV) go build -C apps/bot ./cmd/bot
	$(GOENV) go build -C apps/worker ./cmd/worker
	$(GOENV) go build -C apps/cli ./cmd/cuepoint

test:
	$(GOENV) go test -C apps/api ./...
	$(GOENV) go test -C apps/bot ./...
	$(GOENV) go test -C apps/worker ./...
	$(GOENV) go test -C apps/cli ./...
	$(GOENV) go test -C packages/config ./...
	$(GOENV) go test -C packages/database ./...
	$(GOENV) go test -C packages/events ./...
	$(GOENV) go test -C packages/ffmpeg ./...
	$(GOENV) go test -C packages/stream ./...
