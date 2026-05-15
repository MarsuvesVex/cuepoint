package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarsuvesVex/cuepoint/apps/cli/internal/cli"
	"github.com/MarsuvesVex/cuepoint/packages/config"
	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := cli.NewClient(cfg.CLI.APIBaseURL, nil)
	ctx := context.Background()

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "health":
		if err := client.Health(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("ok")
	case "marker":
		runMarker(ctx, client, os.Args[2:])
	case "job":
		runJob(ctx, client, os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func runMarker(ctx context.Context, client *cli.Client, args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	switch args[0] {
	case "create":
		flags := flag.NewFlagSet("marker create", flag.ExitOnError)
		streamID := flags.String("stream", "", "stream identifier")
		label := flags.String("label", "", "marker label")
		timestamp := flags.String("timestamp", "", "marker timestamp")
		_ = flags.Parse(args[1:])

		result, err := client.CreateMarker(ctx, stream.CreateMarkerInput{
			StreamID:  *streamID,
			Label:     *label,
			Timestamp: *timestamp,
		})
		if err != nil {
			log.Fatal(err)
		}
		printJSON(result)
	case "get":
		flags := flag.NewFlagSet("marker get", flag.ExitOnError)
		id := flags.String("id", "", "marker id")
		_ = flags.Parse(args[1:])
		result, err := client.GetMarker(ctx, *id)
		if err != nil {
			log.Fatal(err)
		}
		printJSON(result)
	default:
		usage()
		os.Exit(2)
	}
}

func runJob(ctx context.Context, client *cli.Client, args []string) {
	if len(args) == 0 || args[0] != "get" {
		usage()
		os.Exit(2)
	}

	flags := flag.NewFlagSet("job get", flag.ExitOnError)
	id := flags.String("id", "", "job id")
	_ = flags.Parse(args[1:])

	result, err := client.GetJob(ctx, *id)
	if err != nil {
		log.Fatal(err)
	}
	printJSON(result)
}

func printJSON(value any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  cuepoint health")
	fmt.Fprintln(os.Stderr, "  cuepoint marker create --stream <id> --label <label> --timestamp <ts>")
	fmt.Fprintln(os.Stderr, "  cuepoint marker get --id <marker-id>")
	fmt.Fprintln(os.Stderr, "  cuepoint job get --id <job-id>")
}
