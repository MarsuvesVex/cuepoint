package bot

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

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

func NewHandler(commands ...Command) *Handler {
	return &Handler{commands: commands}
}

func NewDefaultHandler(markerClient MarkerClient, healthClient HealthcheckClient) *Handler {
	return NewHandler(
		NewHealthcheckCommand(healthClient),
		NewMarkerCommand(markerClient),
	)
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
	messages, errs := adapter.Receive(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return err
			}
		case msg, ok := <-messages:
			if !ok {
				return nil
			}
			reply, err := handler.Handle(ctx, msg)
			if err != nil {
				if respondErr := responder.Reply(ctx, msg, "error: "+err.Error()); respondErr != nil {
					return respondErr
				}
				continue
			}
			if reply == "" {
				continue
			}
			if err := responder.Reply(ctx, msg, reply); err != nil {
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

type HealthcheckCommand struct {
	client HealthcheckClient
}

func NewHealthcheckCommand(client HealthcheckClient) *HealthcheckCommand {
	return &HealthcheckCommand{client: client}
}

func (c *HealthcheckCommand) Name() string {
	return "health"
}

func (c *HealthcheckCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name()
}

func (c *HealthcheckCommand) Run(ctx context.Context, _ Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !health")
	}

	result, err := c.client.Healthcheck(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("health=%s", result.Status), nil
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
