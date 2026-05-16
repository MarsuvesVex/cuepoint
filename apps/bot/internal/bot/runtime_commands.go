package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type SetTitleCommand struct{ client RuntimeClient }
type RestoreTitleCommand struct{ client RuntimeClient }
type ToggleTitlesCommand struct{ client RuntimeClient }
type TitleFormatCommand struct{ client RuntimeClient }
type ViewTitleFormatCommand struct{ client RuntimeClient }
type ResetTitleFormatCommand struct{ client RuntimeClient }
type WatchingCommand struct{ client RuntimeClient }
type ReactCommand struct{ client RuntimeClient }
type NextSegmentCommand struct{ client RuntimeClient }

func NewSetTitleCommand(client RuntimeClient) *SetTitleCommand {
	return &SetTitleCommand{client: client}
}
func NewRestoreTitleCommand(client RuntimeClient) *RestoreTitleCommand {
	return &RestoreTitleCommand{client: client}
}
func NewToggleTitlesCommand(client RuntimeClient) *ToggleTitlesCommand {
	return &ToggleTitlesCommand{client: client}
}
func NewTitleFormatCommand(client RuntimeClient) *TitleFormatCommand {
	return &TitleFormatCommand{client: client}
}
func NewViewTitleFormatCommand(client RuntimeClient) *ViewTitleFormatCommand {
	return &ViewTitleFormatCommand{client: client}
}
func NewResetTitleFormatCommand(client RuntimeClient) *ResetTitleFormatCommand {
	return &ResetTitleFormatCommand{client: client}
}
func NewWatchingCommand(client RuntimeClient) *WatchingCommand {
	return &WatchingCommand{client: client}
}
func NewReactCommand(client RuntimeClient) *ReactCommand { return &ReactCommand{client: client} }
func NewNextSegmentCommand(client RuntimeClient) *NextSegmentCommand {
	return &NextSegmentCommand{client: client}
}

func (c *SetTitleCommand) Name() string         { return "settitle" }
func (c *RestoreTitleCommand) Name() string     { return "restoretitle" }
func (c *ToggleTitlesCommand) Name() string     { return "toggletitles" }
func (c *TitleFormatCommand) Name() string      { return "titleformat" }
func (c *ViewTitleFormatCommand) Name() string  { return "viewtitleformat" }
func (c *ResetTitleFormatCommand) Name() string { return "resettitleformat" }
func (c *WatchingCommand) Name() string         { return "watching" }
func (c *ReactCommand) Name() string            { return "react" }
func (c *NextSegmentCommand) Name() string      { return "nextsegment" }
func (c *SetTitleCommand) Help() string {
	return "!settitle - set the stream title to the current segment title"
}
func (c *RestoreTitleCommand) Help() string {
	return "!restoretitle - restore the original stream title"
}
func (c *ToggleTitlesCommand) Help() string { return "!toggletitles - toggle automatic title updates" }
func (c *TitleFormatCommand) Help() string {
	return "!titleformat {format_string} - set the live title format"
}
func (c *ViewTitleFormatCommand) Help() string {
	return "!viewtitleformat - show the active title format"
}
func (c *ResetTitleFormatCommand) Help() string {
	return "!resettitleformat - reset the live title format"
}
func (c *WatchingCommand) Help() string { return "!watching - show the active react video details" }
func (c *ReactCommand) Help() string    { return "!react - apply the current react segment state" }
func (c *NextSegmentCommand) Help() string {
	return "!nextsegment - end the current segment and advance to the next one"
}
func (c *SetTitleCommand) Match(input ParsedInput) bool     { return input.Name == c.Name() }
func (c *RestoreTitleCommand) Match(input ParsedInput) bool { return input.Name == c.Name() }
func (c *ToggleTitlesCommand) Match(input ParsedInput) bool { return input.Name == c.Name() }
func (c *TitleFormatCommand) Match(input ParsedInput) bool  { return input.Name == c.Name() }
func (c *ViewTitleFormatCommand) Match(input ParsedInput) bool {
	return input.Name == c.Name()
}
func (c *ResetTitleFormatCommand) Match(input ParsedInput) bool { return input.Name == c.Name() }
func (c *WatchingCommand) Match(input ParsedInput) bool         { return input.Name == c.Name() }
func (c *ReactCommand) Match(input ParsedInput) bool            { return input.Name == c.Name() }
func (c *NextSegmentCommand) Match(input ParsedInput) bool      { return input.Name == c.Name() }

func (c *SetTitleCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !settitle")
	}
	state, err := c.client.ApplyCurrentTitle(ctx, messageChannel(msg))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("title=%s", state.Session.CurrentTitle), nil
}

func (c *RestoreTitleCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !restoretitle")
	}
	state, err := c.client.RestoreTitle(ctx, messageChannel(msg))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("title=%s", state.Session.CurrentTitle), nil
}

func (c *ToggleTitlesCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !toggletitles")
	}
	state, err := c.client.ToggleTitles(ctx, messageChannel(msg))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("auto_titles=%t", state.Session.AutoTitleEnabled), nil
}

func (c *TitleFormatCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) == 0 {
		return "", errors.New("usage: !titleformat {format_string}")
	}
	format := strings.Join(input.Args, " ")
	state, err := c.client.SetTitleFormat(ctx, messageChannel(msg), format, false)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("title_format=%s", activeTitleFormat(state)), nil
}

func (c *ViewTitleFormatCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !viewtitleformat")
	}
	result, err := c.client.GetTitleFormat(ctx, messageChannel(msg))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("title_format=%s", result.Format), nil
}

func (c *ResetTitleFormatCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !resettitleformat")
	}
	state, err := c.client.SetTitleFormat(ctx, messageChannel(msg), "", true)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("title_format=%s", activeTitleFormat(state)), nil
}

func (c *WatchingCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !watching")
	}
	state, err := c.client.SyncSession(ctx, messageChannel(msg))
	if err != nil {
		if reply, ok := runtimeReplyForError(err); ok {
			return reply, nil
		}
		return "", err
	}
	if !state.IsLive {
		return "stream=offline", nil
	}
	if state.ActiveSegment == nil || state.ActiveSegment.SegmentType != stream.SegmentTypeReact {
		return "no active react segment", nil
	}
	return fmt.Sprintf("%s | %s - %s", state.ActiveSegment.YouTubeVideoTitle, state.ActiveSegment.YouTubeCreatorName, state.ActiveSegment.YouTubeCanonicalURL), nil
}

func (c *ReactCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !react")
	}
	channel := messageChannel(msg)
	state, err := c.client.SyncSession(ctx, channel)
	if err != nil {
		if reply, ok := runtimeReplyForError(err); ok {
			return reply, nil
		}
		return "", err
	}
	if !state.IsLive {
		return "stream=offline", nil
	}
	if state.ActiveSegment == nil || state.ActiveSegment.SegmentType != stream.SegmentTypeReact {
		return "no active react segment", nil
	}
	state, err = c.client.ApplyCurrentTitle(ctx, channel)
	if err != nil {
		if reply, ok := runtimeReplyForError(err); ok {
			return reply, nil
		}
		return "", err
	}
	return fmt.Sprintf("react=%s title=%s", activeSegmentTitle(state.ActiveSegment), state.Session.CurrentTitle), nil
}

func (c *NextSegmentCommand) Run(ctx context.Context, msg Message, input ParsedInput) (string, error) {
	if len(input.Args) != 0 {
		return "", errors.New("usage: !nextsegment")
	}
	channel := messageChannel(msg)
	state, err := c.client.SyncSession(ctx, channel)
	if err != nil {
		if reply, ok := runtimeReplyForError(err); ok {
			return reply, nil
		}
		return "", err
	}
	if !state.IsLive {
		return "stream=offline", nil
	}
	state, err = c.client.AdvanceSegment(ctx, channel)
	if err != nil {
		if reply, ok := runtimeReplyForError(err); ok {
			return reply, nil
		}
		return "", err
	}
	return fmt.Sprintf("segment=%s", activeSegmentTitle(state.ActiveSegment)), nil
}

func messageChannel(msg Message) string {
	return commandChannel([]string{msg.Channel})
}

func commandChannel(values []string) string {
	if len(values) == 0 {
		return "local"
	}
	channel := strings.TrimPrefix(strings.TrimSpace(values[0]), "#")
	if channel == "" {
		return "local"
	}
	return channel
}

func parseTimelineMarkerArgs(args []string) (label string, end bool) {
	if len(args) > 0 && strings.EqualFold(args[len(args)-1], "[end]") {
		return strings.Join(args[:len(args)-1], " "), true
	}
	return strings.Join(args, " "), false
}

func activeSegmentTitle(segment *stream.PlanSegment) string {
	if segment == nil {
		return ""
	}
	return stream.SegmentDisplayTitle(*segment)
}

func activeTitleFormat(state stream.RuntimeState) string {
	if strings.TrimSpace(state.Session.TitleFormatOverride) != "" {
		return state.Session.TitleFormatOverride
	}
	if strings.TrimSpace(state.Settings.DefaultTitleFormat) != "" {
		return state.Settings.DefaultTitleFormat
	}
	return stream.DefaultTitleFormat()
}

func runtimeReplyForError(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == 404 && strings.EqualFold(strings.TrimSpace(apiErr.Message), "404 Not Found") {
			return "runtime unavailable", true
		}
		switch strings.TrimSpace(apiErr.Message) {
		case "stream is not live":
			return "stream=offline", true
		case "no active segment":
			return "no active segment", true
		case "no next segment":
			return "no next segment", true
		case "not found":
			return "runtime not configured", true
		}
	}
	return "", false
}
