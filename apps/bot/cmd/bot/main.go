package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarsuvesVex/cuepoint/apps/bot/internal/bot"
	"github.com/MarsuvesVex/cuepoint/packages/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := bot.NewHTTPMarkerClient(cfg.Bot.APIBaseURL, nil)
	handler := bot.NewDefaultHandler(client, client)
	adapter := bot.NewStdinAdapter(os.Stdin)
	responder := bot.NewWriterResponder(os.Stdout)

	log.Printf("bot ready, enter commands like !help, !health:all, !health:bot, !health:server, or !marker <stream> <label> <timestamp>")
	if err := bot.Run(ctx, adapter, responder, handler); err != nil {
		log.Fatalf("run bot: %v", err)
	}
}
