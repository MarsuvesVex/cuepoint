module github.com/MarsuvesVex/cuepoint/apps/bot

go 1.26.2

require (
	github.com/MarsuvesVex/cuepoint/packages/config v0.0.0
	github.com/MarsuvesVex/cuepoint/packages/stream v0.0.0
)

replace github.com/MarsuvesVex/cuepoint/packages/config => ../../packages/config

replace github.com/MarsuvesVex/cuepoint/packages/stream => ../../packages/stream
