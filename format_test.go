package slogexp

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestWrapReplaceAttrFunc(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: WrapReplaceAttrFunc(ReplaceTimeAttr(time.TimeOnly), ReplaceSourceAttr()),
		AddSource:   true,
	})
	logger := slog.New(handler)

	logger.Info("info text")
	logger.Warn("warn text")
	logger.Error("error text")
}
