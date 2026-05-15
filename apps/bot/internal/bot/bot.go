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

type CreateMarkerResult struct {
	MarkerID string
	JobID    string
	Status   string
}

type Handler struct {
	client MarkerClient
}

func NewHandler(client MarkerClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Handle(ctx context.Context, msg Message) (string, error) {
	command, ok := parseMarkerCommand(msg.Text)
	if !ok {
		return "", nil
	}

	result, err := h.client.CreateMarker(ctx, command)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("marker=%s job=%s status=%s", result.MarkerID, result.JobID, result.Status), nil
}

func Run(ctx context.Context, adapter Adapter, responder Responder, handler *Handler) error {
	messages, errs := adapter.Receive(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
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

func parseMarkerCommand(text string) (stream.CreateMarkerInput, bool) {
	fields := strings.Fields(text)
	if len(fields) != 4 || fields[0] != "!marker" {
		return stream.CreateMarkerInput{}, false
	}
	return stream.CreateMarkerInput{
		StreamID:  fields[1],
		Label:     fields[2],
		Timestamp: fields[3],
	}, true
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
