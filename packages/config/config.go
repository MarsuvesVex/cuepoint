package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultAPIAddr     = ":8080"
	defaultDatabaseURL = "postgres://cuepoint:cuepoint@127.0.0.1:5432/cuepoint?sslmode=disable"
	defaultRedisAddr   = "127.0.0.1:6379"
	defaultRedisDB     = 0
	defaultWorkerPoll  = 5 * time.Second
	defaultWorkerBlock = 5 * time.Second
	defaultAPIBaseURL  = "http://127.0.0.1:8080"
	defaultBotPrompt   = "> "
	defaultQueueName   = "cuepoint:jobs"
)

type Config struct {
	API       APIConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Worker    WorkerConfig
	Bot       BotConfig
	CLI       CLIConfig
	QueueName string
}

type APIConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Addr string
	DB   int
}

type WorkerConfig struct {
	PollInterval time.Duration
	BlockTimeout time.Duration
}

type BotConfig struct {
	APIBaseURL string
	Prompt     string
}

type CLIConfig struct {
	APIBaseURL string
}

var loadEnvOnce sync.Once

func Load() (Config, error) {
	if err := loadRootEnv(); err != nil {
		return Config{}, err
	}

	redisDB, err := intFromEnv("REDIS_DB", defaultRedisDB)
	if err != nil {
		return Config{}, err
	}

	pollInterval, err := durationFromEnv("WORKER_POLL_INTERVAL", defaultWorkerPoll)
	if err != nil {
		return Config{}, err
	}

	blockTimeout, err := durationFromEnv("WORKER_BLOCK_TIMEOUT", defaultWorkerBlock)
	if err != nil {
		return Config{}, err
	}

	return Config{
		API: APIConfig{
			Addr: stringFromEnv("API_ADDR", defaultAPIAddr),
		},
		Database: DatabaseConfig{
			URL: stringFromEnv("DATABASE_URL", defaultDatabaseURL),
		},
		Redis: RedisConfig{
			Addr: stringFromEnv("REDIS_ADDR", defaultRedisAddr),
			DB:   redisDB,
		},
		Worker: WorkerConfig{
			PollInterval: pollInterval,
			BlockTimeout: blockTimeout,
		},
		Bot: BotConfig{
			APIBaseURL: stringFromEnv("BOT_API_BASE_URL", stringFromEnv("API_BASE_URL", defaultAPIBaseURL)),
			Prompt:     stringFromEnv("BOT_PROMPT", defaultBotPrompt),
		},
		CLI: CLIConfig{
			APIBaseURL: stringFromEnv("CLI_API_BASE_URL", stringFromEnv("API_BASE_URL", defaultAPIBaseURL)),
		},
		QueueName: stringFromEnv("QUEUE_NAME", defaultQueueName),
	}, nil
}

func stringFromEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func intFromEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s: parse int: %w", key, err)
	}

	return parsed, nil
}

func loadRootEnv() error {
	var err error
	loadEnvOnce.Do(func() {
		err = loadEnvFile()
	})
	return err
}

func loadEnvFile() error {
	path, err := findEnvFile(".env")
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid env line %q", line)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("invalid env key in line %q", line)
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s from %s: %w", key, path, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %s: %w", path, err)
	}

	return nil
}

func findEnvFile(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

func durationFromEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s: parse duration: %w", key, err)
	}

	return parsed, nil
}
