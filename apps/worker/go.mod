module github.com/MarsuvesVex/cuepoint/apps/worker

go 1.26.2

require (
	github.com/MarsuvesVex/cuepoint/packages/config v0.0.0
	github.com/MarsuvesVex/cuepoint/packages/database v0.0.0
	github.com/MarsuvesVex/cuepoint/packages/events v0.0.0
	github.com/MarsuvesVex/cuepoint/packages/ffmpeg v0.0.0
	github.com/MarsuvesVex/cuepoint/packages/stream v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/MarsuvesVex/cuepoint/packages/config => ../../packages/config

replace github.com/MarsuvesVex/cuepoint/packages/database => ../../packages/database

replace github.com/MarsuvesVex/cuepoint/packages/events => ../../packages/events

replace github.com/MarsuvesVex/cuepoint/packages/ffmpeg => ../../packages/ffmpeg

replace github.com/MarsuvesVex/cuepoint/packages/stream => ../../packages/stream
