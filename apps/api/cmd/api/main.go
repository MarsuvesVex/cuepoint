package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarsuvesVex/cuepoint/apps/api/internal/api"
	"github.com/MarsuvesVex/cuepoint/packages/config"
	"github.com/MarsuvesVex/cuepoint/packages/database"
	"github.com/MarsuvesVex/cuepoint/packages/events"
	"github.com/MarsuvesVex/cuepoint/packages/stream"
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

	service := stream.NewService(store, queue)
	runtimeService := api.NewRuntimeService(
		store,
		api.NewTwitchBridgeClient(cfg.Web.BaseURL, cfg.Internal.HeaderName, cfg.Internal.ServiceToken, nil),
	)
	server := &http.Server{
		Addr:    cfg.API.Addr,
		Handler: api.NewServer(service, store, runtimeService, cfg.Internal.HeaderName, cfg.Internal.ServiceToken).Handler(),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("api listening on %s", cfg.API.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
