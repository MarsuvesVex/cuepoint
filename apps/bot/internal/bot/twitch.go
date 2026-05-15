package bot

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/MarsuvesVex/cuepoint/packages/config"
)

type TwitchAdapter struct {
	cfg    config.TwitchConfig
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	logger *Logger
	mu     sync.Mutex
}

func NewTwitchAdapter(ctx context.Context, cfg config.TwitchConfig) (*TwitchAdapter, error) {
	return NewTwitchAdapterWithLogger(ctx, cfg, nil)
}

func NewTwitchAdapterWithLogger(ctx context.Context, cfg config.TwitchConfig, logger *Logger) (*TwitchAdapter, error) {
	if logger == nil {
		logger = NewLogger("info", nil)
	}
	if cfg.Username == "" {
		return nil, errors.New("BOT_TWITCH_USERNAME is required")
	}
	if cfg.OAuthToken == "" {
		return nil, errors.New("BOT_TWITCH_OAUTH_TOKEN is required")
	}
	if cfg.Channel == "" {
		return nil, errors.New("BOT_TWITCH_CHANNEL is required")
	}
	if strings.Contains(cfg.Channel, ":") {
		return nil, errors.New("BOT_TWITCH_CHANNEL looks like an address; did you mean to set BOT_TWITCH_ADDR instead")
	}

	var (
		conn net.Conn
		err  error
	)
	dialer := &net.Dialer{}
	logger.Infof("connecting to twitch addr=%s channel=%s tls=%t user=%s", cfg.Addr, normalizeChannel(cfg.Channel), cfg.UseTLS, cfg.Username)
	if cfg.UseTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", cfg.Addr, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: hostPart(cfg.Addr),
		})
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", cfg.Addr)
	}
	if err != nil {
		return nil, fmt.Errorf("connect to twitch irc: %w", err)
	}
	logger.Infof("connected to twitch irc addr=%s", cfg.Addr)

	a := &TwitchAdapter{
		cfg:    cfg,
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		logger: logger,
	}

	if err := a.writeLine("PASS " + normalizeOAuthToken(cfg.OAuthToken)); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := a.writeLine("NICK " + cfg.Username); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := a.writeLine("JOIN #" + normalizeChannel(cfg.Channel)); err != nil {
		_ = conn.Close()
		return nil, err
	}
	logger.Infof("joined twitch channel=%s", normalizeChannel(cfg.Channel))

	return a, nil
}

func (a *TwitchAdapter) Receive(ctx context.Context) (<-chan Message, <-chan error) {
	msgs := make(chan Message)
	errs := make(chan error, 1)

	go func() {
		defer close(msgs)
		defer close(errs)
		defer a.Close()

		go func() {
			<-ctx.Done()
			_ = a.Close()
		}()

		for {
			line, err := a.reader.ReadString('\n')
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				errs <- fmt.Errorf("read twitch irc: %w", err)
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if strings.HasPrefix(line, "PING ") {
				a.logger.Debugf("received twitch ping line=%q", line)
				if err := a.writeLine("PONG " + strings.TrimPrefix(line, "PING ")); err != nil {
					errs <- err
					return
				}
				continue
			}

			msg, ok := parseTwitchPrivmsg(line)
			if !ok {
				a.logger.Debugf("ignoring twitch irc line=%q", line)
				continue
			}
			a.logger.Debugf("parsed twitch privmsg channel=%s user=%s text=%q", msg.Channel, msg.User, msg.Text)

			select {
			case <-ctx.Done():
				return
			case msgs <- msg:
			}
		}
	}()

	return msgs, errs
}

func (a *TwitchAdapter) Reply(_ context.Context, msg Message, text string) error {
	channel := msg.Channel
	if channel == "" {
		channel = normalizeChannel(a.cfg.Channel)
	}
	a.logger.Debugf("replying to twitch channel=%s text=%q", normalizeChannel(channel), text)
	return a.writeLine(fmt.Sprintf("PRIVMSG #%s :%s", normalizeChannel(channel), text))
}

func (a *TwitchAdapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn == nil {
		return nil
	}
	a.logger.Infof("closing twitch connection")
	err := a.conn.Close()
	a.conn = nil
	return err
}

func (a *TwitchAdapter) writeLine(line string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn == nil {
		return errors.New("twitch connection is closed")
	}
	if _, err := a.writer.WriteString(line + "\r\n"); err != nil {
		return fmt.Errorf("write twitch irc line: %w", err)
	}
	if err := a.writer.Flush(); err != nil {
		return fmt.Errorf("flush twitch irc line: %w", err)
	}
	a.logger.Debugf("sent twitch irc line=%q", line)
	return nil
}

func parseTwitchPrivmsg(line string) (Message, bool) {
	prefix, payload, ok := strings.Cut(line, " PRIVMSG ")
	if !ok || !strings.HasPrefix(prefix, ":") {
		return Message{}, false
	}

	userPrefix := strings.TrimPrefix(prefix, ":")
	user, _, ok := strings.Cut(userPrefix, "!")
	if !ok {
		return Message{}, false
	}

	channelPart, text, ok := strings.Cut(payload, " :")
	if !ok {
		return Message{}, false
	}

	channel := strings.TrimPrefix(strings.TrimSpace(channelPart), "#")
	return Message{
		Channel: channel,
		User:    user,
		Text:    text,
	}, true
}

func normalizeOAuthToken(token string) string {
	if strings.HasPrefix(token, "oauth:") {
		return token
	}
	return "oauth:" + token
}

func normalizeChannel(channel string) string {
	return strings.TrimPrefix(strings.TrimSpace(channel), "#")
}

func hostPart(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
