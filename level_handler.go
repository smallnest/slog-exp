package slogexp

import (
	"context"
	"log/slog"
)

// LevelHandler is a handler that logs records with levels greater than or
// equal to a given level using the given handler.
type LevelHandler struct {
	level         slog.Level // minimum level to log
	levelHandlers map[slog.Level]slog.Handler
}

// NewLevelHandler returns a new handler that logs records with levels greater
// than or equal to the given level using the given handler.
func NewLevelHandler(levelHandlers map[slog.Level]slog.Handler) *LevelHandler {
	var level = slog.LevelDebug

	for l := range levelHandlers {
		if l < level {
			level = l
		}
	}

	return &LevelHandler{
		level:         level,
		levelHandlers: levelHandlers,
	}
}

func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle logs the record if the level is greater than or equal to the
// handler's level.
func (h *LevelHandler) Handle(ctx context.Context, record slog.Record) error {
	if handler, ok := h.levelHandlers[record.Level]; ok && handler != nil {
		return handler.Handle(ctx, record)
	}
	return nil
}

// WithAttrs returns a new handler with the given attributes added to the
// existing attributes.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	levelHandlers := make(map[slog.Level]slog.Handler, len(h.levelHandlers))
	for k, v := range h.levelHandlers {
		levelHandlers[k] = v.WithAttrs(attrs)
	}
	return &LevelHandler{
		level:         h.level,
		levelHandlers: levelHandlers,
	}
}

// WithGroup returns a new handler with the given group added to the existing
// group.
func (h *LevelHandler) WithGroup(name string) slog.Handler {
	levelHandlers := make(map[slog.Level]slog.Handler, len(h.levelHandlers))
	for k, v := range h.levelHandlers {
		levelHandlers[k] = v.WithGroup(name)
	}
	return &LevelHandler{
		level:         h.level,
		levelHandlers: levelHandlers,
	}
}
