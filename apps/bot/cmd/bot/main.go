package main

import (
	"context"
	"fmt"
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

	botLogger := bot.NewLogger(cfg.Bot.LogLevel, log.Default())
	client := bot.NewHTTPMarkerClient(cfg.Bot.APIBaseURL, nil)
	handler := bot.NewDefaultHandler(client, client)
	adapter, responder, mode, err := buildTransport(ctx, cfg, botLogger)
	if err != nil {
		log.Fatalf("build bot transport: %v", err)
	}

	botLogger.Infof("bot ready with transport=%s api_base_url=%s log_level=%s", mode, cfg.Bot.APIBaseURL, cfg.Bot.LogLevel)
	if err := bot.RunWithLogger(ctx, adapter, responder, handler, botLogger); err != nil {
		log.Fatalf("run bot: %v", err)
	}
}

func buildTransport(ctx context.Context, cfg config.Config, logger *bot.Logger) (bot.Adapter, bot.Responder, string, error) {
	switch cfg.Bot.Transport {
	case "stdin", "":
		return bot.NewStdinAdapter(os.Stdin), bot.NewWriterResponder(os.Stdout), "stdin", nil
	case "twitch":
		adapter, err := bot.NewTwitchAdapterWithLogger(ctx, cfg.Bot.Twitch, logger)
		if err != nil {
			return nil, nil, "", err
		}
		return adapter, adapter, "twitch", nil
	default:
		return nil, nil, "", fmt.Errorf("unsupported BOT_TRANSPORT %q", cfg.Bot.Transport)
	}
}
