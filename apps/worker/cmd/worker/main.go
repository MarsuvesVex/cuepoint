package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarsuvesVex/cuepoint/apps/worker/internal/worker"
	"github.com/MarsuvesVex/cuepoint/packages/config"
	"github.com/MarsuvesVex/cuepoint/packages/database"
	"github.com/MarsuvesVex/cuepoint/packages/events"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	store, err := database.Open(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer store.Close()

	if err := store.Bootstrap(ctx); err != nil {
		log.Fatalf("bootstrap database: %v", err)
	}

	redisClient := events.NewRedisClient(cfg.Redis.Addr, cfg.Redis.DB)
	defer func() { _ = redisClient.Close() }()
	queue := events.NewQueue(redisClient, cfg.QueueName)

	processor := worker.NewProcessor(store, queue, cfg.Worker.BlockTimeout, log.Default())
	if err := processor.Run(ctx); err != nil {
		log.Fatalf("run worker: %v", err)
	}
}
