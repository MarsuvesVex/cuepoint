package bot

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type fakeClient struct {
	createResult CreateMarkerResult
	createErr    error
	healthResult HealthcheckResult
	healthErr    error
}

func (c fakeClient) CreateMarker(context.Context, stream.CreateMarkerInput) (CreateMarkerResult, error) {
	return c.createResult, c.createErr
}

func (c fakeClient) Healthcheck(context.Context) (HealthcheckResult, error) {
	return c.healthResult, c.healthErr
}

func TestParseCommand(t *testing.T) {
	input, ok := ParseCommand("!marker stream clip 00:00:10")
	if !ok {
		t.Fatal("expected command to parse")
	}
	if input.Name != "marker" {
		t.Fatalf("Name = %q", input.Name)
	}
	if len(input.Args) != 3 {
		t.Fatalf("Args = %#v", input.Args)
	}
}

func TestParseCommandRejectsNonCommands(t *testing.T) {
	if _, ok := ParseCommand("hello world"); ok {
		t.Fatal("expected non-command to be rejected")
	}
}

func TestMarkerCommand(t *testing.T) {
	handler := NewDefaultHandler(fakeClient{
		createResult: CreateMarkerResult{
			MarkerID: "marker-1",
			JobID:    "job-1",
			Status:   "pending",
		},
	}, fakeClient{})

	reply, err := handler.Handle(context.Background(), Message{Text: "!marker stream clip 00:00:10"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "marker=marker-1 job=job-1 status=pending" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthCommand(t *testing.T) {
	client := fakeClient{healthResult: HealthcheckResult{Status: "ok"}}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health:all"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "bot=ok server=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthBotCommand(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health:bot"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "bot=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthServerCommand(t *testing.T) {
	client := fakeClient{healthResult: HealthcheckResult{Status: "ok"}}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health:server"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "server=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthServerTypoAlias(t *testing.T) {
	client := fakeClient{healthResult: HealthcheckResult{Status: "ok"}}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!heath:server"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "server=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHelpCommand(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!help"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	for _, want := range []string{"!help", "!health:all", "!health:bot", "!health:server", "!marker"} {
		if !strings.Contains(reply, want) {
			t.Fatalf("reply %q missing %q", reply, want)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!nope"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "unknown command: nope" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestMarkerUsageError(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client)

	_, err := handler.Handle(context.Background(), Message{Text: "!marker stream clip"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHealthError(t *testing.T) {
	client := fakeClient{healthErr: errors.New("boom")}
	handler := NewDefaultHandler(client, client)

	_, err := handler.Handle(context.Background(), Message{Text: "!health:all"})
	if err == nil {
		t.Fatal("expected error")
	}
}
