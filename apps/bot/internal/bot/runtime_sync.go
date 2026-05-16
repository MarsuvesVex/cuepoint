package bot

import (
	"context"
	"errors"
	"strings"
	"time"
)

func RunRuntimeSyncLoop(ctx context.Context, client RuntimeClient, channel string, interval time.Duration, logger *Logger) {
	if client == nil || strings.TrimSpace(channel) == "" || interval <= 0 {
		return
	}
	if logger == nil {
		logger = NewLogger("info", nil)
	}

	channel = strings.TrimPrefix(strings.TrimSpace(channel), "#")
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := client.SyncSession(ctx, channel); err != nil {
					var apiErr *APIError
					if errors.As(err, &apiErr) {
						if apiErr.StatusCode == 404 {
							logger.Debugf("runtime sync unavailable for channel=%s: %v", channel, err)
							continue
						}
						if reply, ok := runtimeReplyForError(err); ok && (reply == "runtime not configured" || reply == "stream=offline") {
							logger.Debugf("runtime sync skipped for channel=%s: %s", channel, reply)
							continue
						}
					}
					logger.Warnf("runtime sync failed for channel=%s: %v", channel, err)
				}
			}
		}
	}()
}
