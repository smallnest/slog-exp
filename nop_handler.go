package slogexp

import (
	"context"
	"log/slog"
)

// NewNopHandler creates a slog handler that discards all log messages.
func NewNopHandler() slog.Handler {
	return &nopHandler{}
}

var _ slog.Handler = (*nopHandler)(nil)

type nopHandler struct{}

func (n *nopHandler) Enabled(_ context.Context, _ slog.Level) bool  { return false }
func (n *nopHandler) Handle(_ context.Context, _ slog.Record) error { return nil }
func (n *nopHandler) WithAttrs(_ []slog.Attr) slog.Handler          { return n }
func (n *nopHandler) WithGroup(_ string) slog.Handler               { return n }

// NewNopLogger creates a slog logger that discards all log messages.
func NewNopLogger() *slog.Logger {
	return slog.New(NewNopHandler())
}
