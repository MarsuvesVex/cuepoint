package bot

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type Message struct {
	Channel string
	User    string
	Text    string
}

type Adapter interface {
	Receive(ctx context.Context) (<-chan Message, <-chan error)
}

type Responder interface {
	Reply(ctx context.Context, msg Message, text string) error
}

type MarkerClient interface {
	CreateMarker(ctx context.Context, input stream.CreateMarkerInput) (CreateMarkerResult, error)
}

type HealthcheckClient interface {
	Healthcheck(ctx context.Context) (HealthcheckResult, error)
}

type CreateMarkerResult struct {
	MarkerID string
	JobID    string
	Status   string
}

type Command interface {
	Name() string
	Help() string
	Match(input ParsedInput) bool
	Run(ctx context.Context, msg Message, input ParsedInput) (string, error)
}

type ParsedInput struct {
	Name string
	Args []string
}

type Handler struct {
	commands []Command
}

type RuntimeStatus struct {
	startedAt time.Time
	now       func() time.Time
}

func NewHandler(commands ...Command) *Handler {
	return &Handler{commands: commands}
}

func NewDefaultHandler(markerClient MarkerClient, healthClient HealthcheckClient) *Handler {
	runtime := NewRuntimeStatus()
	commands := []Command{
		NewHealthAllCommand(healthClient, runtime),
		NewHealthBotCommand(runtime),
		NewHealthServerCommand(healthClient),
		NewMarkerCommand(markerClient),
	}
	return NewHandler(append([]Command{NewHelpCommand(commands)}, commands...)...)
}

func NewRuntimeStatus() *RuntimeStatus {
	return &RuntimeStatus{
		startedAt: time.Now().UTC(),
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (r *RuntimeStatus) Uptime() time.Duration {
	return r.now().Sub(r.startedAt).Round(time.Second)
}

func (h *Handler) Handle(ctx context.Context, msg Message) (string, error) {
	input, ok := ParseCommand(msg.Text)
	if !ok {
		return "", nil
	}

	for _, command := range h.commands {
		if command.Match(input) {
			return command.Run(ctx, msg, input)
		}
	}

	return fmt.Sprintf("unknown command: %s", input.Name), nil
}

func Run(ctx context.Context, adapter Adapter, responder Responder, handler *Handler) error {
	return RunWithLogger(ctx, adapter, responder, handler, nil)
}

func RunWithLogger(ctx context.Context, adapter Adapter, responder Responder, handler *Handler, logger *Logger) error {
	if logger == nil {
		logger = NewLogger("info", nil)
	}

	messages, errs := adapter.Receive(ctx)
	for {
		select {
		case <-ctx.Done():
			logger.Infof("shutting down bot loop")
			return nil
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				logger.Errorf("adapter receive error: %v", err)
				return err
			}
		case msg, ok := <-messages:
			if !ok {
				logger.Infof("adapter message channel closed")
				return nil
			}
			logger.Debugf("received message channel=%s user=%s text=%q", msg.Channel, msg.User, msg.Text)
			reply, err := handler.Handle(ctx, msg)
			if err != nil {
				logger.Warnf("command handling failed for %q: %v", msg.Text, err)
				if respondErr := responder.Reply(ctx, msg, "error: "+err.Error()); respondErr != nil {
					logger.Errorf("failed to send error reply: %v", respondErr)
					return respondErr
				}
				continue
			}
			if reply == "" {
				logger.Debugf("no reply generated for %q", msg.Text)
				continue
			}
			logger.Debugf("sending reply channel=%s user=%s text=%q", msg.Channel, msg.User, reply)
			if err := responder.Reply(ctx, msg, reply); err != nil {
				logger.Errorf("failed to send reply: %v", err)
				return err
			}
		}
	}
}

func ParseCommand(text string) (ParsedInput, bool) {
	fields := strings.Fields(text)
	if len(fields) == 0 || !strings.HasPrefix(fields[0], "!") {
		return ParsedInput{}, false
	}

	return ParsedInput{
		Name: strings.TrimPrefix(fields[0], "!"),
		Args: fields[1:],
	}, true
}

type MarkerCommand struct {
	client MarkerClient
}

func NewMarkerCommand(client MarkerClient) *MarkerCommand {
	return &MarkerCommand{client: client}
}

func (c *MarkerCommand) Name() string {
	return "marker"
}

func (c *MarkerCommand) Help() string {
	return "!marker <stream> <label> <timestamp> - create a marker and queue a job"
}

func (c *MarkerCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name()
}

func (c *MarkerCommand) Run(ctx context.Context, _ Message, input ParsedInput) (string, error) {
	command, err := parseMarkerArgs(input.Args)
	if err != nil {
		return "", err
	}

	result, err := c.client.CreateMarker(ctx, command)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("marker=%s job=%s status=%s", result.MarkerID, result.JobID, result.Status), nil
}

type HealthAllCommand struct {
	client  HealthcheckClient
	runtime *RuntimeStatus
}

func NewHealthAllCommand(client HealthcheckClient, runtime *RuntimeStatus) *HealthAllCommand {
	return &HealthAllCommand{client: client, runtime: runtime}
}

func (c *HealthAllCommand) Name() string {
	return "health:all"
}

func (c *HealthAllCommand) Help() string {
	return "!health:all - show bot and server health"
}

func (c *HealthAllCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name() || input.Name == "health"
}

func (c *HealthAllCommand) Run(ctx context.Context, _ Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !health:all")
	}

	result, err := c.client.Healthcheck(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("bot=ok uptime=%s server=%s", c.runtime.Uptime(), result.Status), nil
}

type HealthBotCommand struct {
	runtime *RuntimeStatus
}

func NewHealthBotCommand(runtime *RuntimeStatus) *HealthBotCommand {
	return &HealthBotCommand{runtime: runtime}
}

func (c *HealthBotCommand) Name() string {
	return "health:bot"
}

func (c *HealthBotCommand) Help() string {
	return "!health:bot - show bot health"
}

func (c *HealthBotCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name()
}

func (c *HealthBotCommand) Run(_ context.Context, _ Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !health:bot")
	}
	return fmt.Sprintf("bot=ok uptime=%s", c.runtime.Uptime()), nil
}

type HealthServerCommand struct {
	client HealthcheckClient
}

func NewHealthServerCommand(client HealthcheckClient) *HealthServerCommand {
	return &HealthServerCommand{client: client}
}

func (c *HealthServerCommand) Name() string {
	return "health:server"
}

func (c *HealthServerCommand) Help() string {
	return "!health:server - show API server health"
}

func (c *HealthServerCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name() || input.Name == "heath:server"
}

func (c *HealthServerCommand) Run(ctx context.Context, _ Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !health:server")
	}

	result, err := c.client.Healthcheck(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("server=%s", result.Status), nil
}

type HelpCommand struct {
	commands []Command
}

func NewHelpCommand(commands []Command) *HelpCommand {
	return &HelpCommand{commands: commands}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Help() string {
	return "!help - show available commands"
}

func (c *HelpCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name()
}

func (c *HelpCommand) Run(_ context.Context, _ Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !help")
	}

	lines := []string{c.Help()}
	for _, command := range c.commands {
		lines = append(lines, command.Help())
	}
	return strings.Join(lines, " | "), nil
}

func parseMarkerArgs(args []string) (stream.CreateMarkerInput, error) {
	if len(args) != 3 {
		return stream.CreateMarkerInput{}, errors.New("usage: !marker <stream> <label> <timestamp>")
	}
	return stream.CreateMarkerInput{
		StreamID:  args[0],
		Label:     args[1],
		Timestamp: args[2],
	}, nil
}

type StdinAdapter struct {
	reader io.Reader
}

func NewStdinAdapter(reader io.Reader) *StdinAdapter {
	return &StdinAdapter{reader: reader}
}

func (a *StdinAdapter) Receive(ctx context.Context) (<-chan Message, <-chan error) {
	msgs := make(chan Message)
	errs := make(chan error, 1)

	go func() {
		defer close(msgs)
		defer close(errs)

		scanner := bufio.NewScanner(a.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case msgs <- Message{Channel: "local", User: "local", Text: scanner.Text()}:
			}
		}

		if err := scanner.Err(); err != nil {
			errs <- err
		}
	}()

	return msgs, errs
}

type WriterResponder struct {
	writer io.Writer
}

func NewWriterResponder(writer io.Writer) *WriterResponder {
	return &WriterResponder{writer: writer}
}

func (r *WriterResponder) Reply(_ context.Context, _ Message, text string) error {
	if r.writer == nil {
		return errors.New("writer is nil")
	}
	_, err := fmt.Fprintln(r.writer, text)
	return err
}
