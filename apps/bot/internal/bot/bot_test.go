package bot

import (
	"bufio"
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/config"
	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type fakeClient struct {
	createResult CreateMarkerResult
	createErr    error
	healthResult HealthcheckResult
	healthErr    error
	runtimeState stream.RuntimeState
	titleFormat  RuntimeTitleFormatResult
	runtimeErr   error
}

func (c fakeClient) CreateMarker(context.Context, stream.CreateMarkerInput) (CreateMarkerResult, error) {
	return c.createResult, c.createErr
}

func (c fakeClient) Healthcheck(context.Context) (HealthcheckResult, error) {
	return c.healthResult, c.healthErr
}

func (c fakeClient) SyncSession(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) GetRuntime(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) ApplyCurrentTitle(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) RestoreTitle(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) ToggleTitles(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) SetTitleFormat(context.Context, string, string, bool) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) GetTitleFormat(context.Context, string) (RuntimeTitleFormatResult, error) {
	return c.titleFormat, nil
}

func (c fakeClient) AdvanceSegment(context.Context, string) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
}

func (c fakeClient) AddTimelineMarker(context.Context, string, string, bool) (stream.RuntimeState, error) {
	return c.runtimeState, c.runtimeErr
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
	}, fakeClient{}, fakeClient{})

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
	runtime := &RuntimeStatus{
		startedAt: time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC),
		now: func() time.Time {
			return time.Date(2026, 5, 16, 0, 0, 45, 0, time.UTC)
		},
	}
	handler := NewHandler(
		NewHelpCommand([]Command{
			NewHealthAllCommand(client, runtime),
			NewHealthBotCommand(runtime),
			NewHealthServerCommand(client),
			NewMarkerCommand(client, client),
		}),
		NewHealthAllCommand(client, runtime),
		NewHealthBotCommand(runtime),
		NewHealthServerCommand(client),
		NewMarkerCommand(client, client),
	)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health:all"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "bot=ok uptime=45s server=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthBotCommand(t *testing.T) {
	client := fakeClient{}
	runtime := &RuntimeStatus{
		startedAt: time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC),
		now: func() time.Time {
			return time.Date(2026, 5, 16, 0, 1, 5, 0, time.UTC)
		},
	}
	handler := NewHandler(
		NewHelpCommand([]Command{
			NewHealthAllCommand(client, runtime),
			NewHealthBotCommand(runtime),
			NewHealthServerCommand(client),
			NewMarkerCommand(client, client),
		}),
		NewHealthAllCommand(client, runtime),
		NewHealthBotCommand(runtime),
		NewHealthServerCommand(client),
		NewMarkerCommand(client, client),
	)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health:bot"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "bot=ok uptime=1m5s" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthAliasCommand(t *testing.T) {
	client := fakeClient{healthResult: HealthcheckResult{Status: "ok"}}
	runtime := &RuntimeStatus{
		startedAt: time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC),
		now: func() time.Time {
			return time.Date(2026, 5, 16, 0, 0, 9, 0, time.UTC)
		},
	}
	handler := NewHandler(
		NewHelpCommand([]Command{
			NewHealthAllCommand(client, runtime),
			NewHealthBotCommand(runtime),
			NewHealthServerCommand(client),
			NewMarkerCommand(client, client),
		}),
		NewHealthAllCommand(client, runtime),
		NewHealthBotCommand(runtime),
		NewHealthServerCommand(client),
		NewMarkerCommand(client, client),
	)

	reply, err := handler.Handle(context.Background(), Message{Text: "!health"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "bot=ok uptime=9s server=ok" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestHealthServerCommand(t *testing.T) {
	client := fakeClient{healthResult: HealthcheckResult{Status: "ok"}}
	handler := NewDefaultHandler(client, client, client)

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
	handler := NewDefaultHandler(client, client, client)

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
	handler := NewDefaultHandler(client, client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!help"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	for _, want := range []string{"!help", "!health", "!marker", "!react", "!nextsegment"} {
		if !strings.Contains(reply, want) {
			t.Fatalf("reply %q missing %q", reply, want)
		}
	}
	if strings.Contains(reply, "\n") {
		t.Fatalf("reply %q should be single-line", reply)
	}
}

func TestHelpTopicCommand(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client, client)

	reply, err := handler.Handle(context.Background(), Message{Text: "!help runtime"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	for _, want := range []string{"!react", "!watching", "!nextsegment"} {
		if !strings.Contains(reply, want) {
			t.Fatalf("reply %q missing %q", reply, want)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	client := fakeClient{}
	handler := NewDefaultHandler(client, client, client)

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
	handler := NewDefaultHandler(client, client, client)

	_, err := handler.Handle(context.Background(), Message{Text: "!marker"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHealthError(t *testing.T) {
	client := fakeClient{healthErr: errors.New("boom")}
	handler := NewDefaultHandler(client, client, client)

	_, err := handler.Handle(context.Background(), Message{Text: "!health:all"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNextSegmentOfflineReply(t *testing.T) {
	client := fakeClient{
		runtimeErr: &APIError{StatusCode: 400, Message: "stream is not live"},
	}
	handler := NewDefaultHandler(client, client, client)

	reply, err := handler.Handle(context.Background(), Message{Channel: "cuepoint", Text: "!nextsegment"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "stream=offline" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestNextSegmentRuntimeUnavailableReply(t *testing.T) {
	client := fakeClient{
		runtimeErr: &APIError{StatusCode: 404, Message: "404 Not Found"},
	}
	handler := NewDefaultHandler(client, client, client)

	reply, err := handler.Handle(context.Background(), Message{Channel: "cuepoint", Text: "!nextsegment"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "runtime unavailable" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestReactNoActiveReactSegmentReply(t *testing.T) {
	client := fakeClient{
		runtimeState: stream.RuntimeState{
			IsLive: true,
			ActiveSegment: &stream.PlanSegment{
				SegmentType:  stream.SegmentTypeStandard,
				SegmentTitle: "Intro",
			},
		},
	}
	handler := NewDefaultHandler(client, client, client)

	reply, err := handler.Handle(context.Background(), Message{Channel: "cuepoint", Text: "!react"})
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if reply != "no active react segment" {
		t.Fatalf("reply = %q", reply)
	}
}

func TestParseTwitchPrivmsg(t *testing.T) {
	msg, ok := parseTwitchPrivmsg(":testuser!testuser@testuser.tmi.twitch.tv PRIVMSG #cuepoint :!help")
	if !ok {
		t.Fatal("expected message to parse")
	}
	if msg.User != "testuser" || msg.Channel != "cuepoint" || msg.Text != "!help" {
		t.Fatalf("unexpected message: %+v", msg)
	}
}

func TestNormalizeOAuthToken(t *testing.T) {
	if got := normalizeOAuthToken("token"); got != "oauth:token" {
		t.Fatalf("token = %q", got)
	}
	if got := normalizeOAuthToken("oauth:token"); got != "oauth:token" {
		t.Fatalf("token = %q", got)
	}
}

func TestTwitchReply(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	adapter := &TwitchAdapter{
		cfg:    config.TwitchConfig{Channel: "cuepoint"},
		conn:   client,
		reader: bufio.NewReader(client),
		writer: bufio.NewWriter(client),
	}
	defer adapter.Close()

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 256)
		n, _ := server.Read(buf)
		done <- string(buf[:n])
	}()

	if err := adapter.Reply(context.Background(), Message{Channel: "cuepoint"}, "hello"); err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}

	if line := <-done; line != "PRIVMSG #cuepoint :hello\r\n" {
		t.Fatalf("line = %q", line)
	}
}

func TestNewTwitchAdapterRejectsAddressAsChannel(t *testing.T) {
	_, err := NewTwitchAdapter(context.Background(), config.TwitchConfig{
		Username:   "botuser",
		OAuthToken: "oauth:token",
		Channel:    "irc.chat.twitch.tv:6697",
		Addr:       "irc.chat.twitch.tv:6697",
		UseTLS:     true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "did you mean to set BOT_TWITCH_ADDR") {
		t.Fatalf("err = %v", err)
	}
}
