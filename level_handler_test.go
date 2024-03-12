package slogexp

import (
	"bytes"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/smallnest/gid"
)

func ExampleLevelHandler() {
	infoFile, err := os.OpenFile("testdata/info.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer infoFile.Close()

	warnFile, err := os.OpenFile("testdata/warn.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer warnFile.Close()

	errorFile, err := os.OpenFile("testdata/error.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer errorFile.Close()

	infoHandler := slog.NewTextHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo})
	warnHandler := slog.NewTextHandler(warnFile, &slog.HandlerOptions{Level: slog.LevelWarn})
	errorHandler := slog.NewTextHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError})

	handler := NewLevelHandler(map[slog.Level]slog.Handler{
		slog.LevelInfo:  infoHandler,
		slog.LevelWarn:  warnHandler,
		slog.LevelError: errorHandler,
	})

	logger := slog.New(handler)

	logger.Info("info text")
	logger.Warn("warn text")
	logger.Error("error text")

	// Output:
}

func TestLevelHandler(t *testing.T) {
	var infoLog bytes.Buffer
	var warnLog bytes.Buffer
	var errorLog bytes.Buffer

	infoHandler := slog.NewTextHandler(&infoLog, &slog.HandlerOptions{Level: slog.LevelInfo})
	warnHandler := slog.NewTextHandler(&warnLog, &slog.HandlerOptions{Level: slog.LevelWarn})
	errorHandler := slog.NewTextHandler(&errorLog, &slog.HandlerOptions{Level: slog.LevelError})

	handler := NewLevelHandler(map[slog.Level]slog.Handler{
		slog.LevelInfo:  infoHandler,
		slog.LevelWarn:  warnHandler,
		slog.LevelError: errorHandler,
	})

	logger := slog.New(handler)

	logger.Info("info text")
	logger.Warn("warn text")
	logger.Error("error text")

	if !strings.Contains(infoLog.String(), `level=INFO msg="info text"`) {
		t.Errorf("infoLog: expected 'level=INFO msg=\"info text\"', got '%s'", infoLog.String())
	}

	if !strings.Contains(warnLog.String(), `level=WARN msg="warn text"`) {
		t.Errorf("warnLog: expected 'level=WARN msg=\"warn text\"', got '%s'", warnLog.String())
	}

	if !strings.Contains(errorLog.String(), `level=ERROR msg="error text"`) {
		t.Errorf("errorLog: expected 'level=ERROR msg=\"error text\"', got '%s'", errorLog.String())
	}
}

func TestLevelHandler_WithGroup(t *testing.T) {
	infoFile, err := os.OpenFile("testdata/info.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer infoFile.Close()

	warnFile, err := os.OpenFile("testdata/warn.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer warnFile.Close()

	errorFile, err := os.OpenFile("testdata/error.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer errorFile.Close()

	infoHandler := slog.NewTextHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo})
	warnHandler := slog.NewTextHandler(warnFile, &slog.HandlerOptions{Level: slog.LevelWarn})
	errorHandler := slog.NewTextHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError})

	handler := NewLevelHandler(map[slog.Level]slog.Handler{
		slog.LevelInfo:  infoHandler,
		slog.LevelWarn:  warnHandler,
		slog.LevelError: errorHandler,
	})

	buildInfo, _ := debug.ReadBuildInfo()

	groupLogger := slog.New(handler).WithGroup("program_info")
	childLogger := groupLogger.With(
		slog.Int("pid", os.Getpid()),
		slog.String("go_version", buildInfo.GoVersion),
	)

	childLogger.Info("info text")
	childLogger.Warn("warn text")
	childLogger.Error("error text")
}

func TestLevelHandler_WithAttr(t *testing.T) {
	infoFile, err := os.OpenFile("testdata/info.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer infoFile.Close()

	warnFile, err := os.OpenFile("testdata/warn.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer warnFile.Close()

	errorFile, err := os.OpenFile("testdata/error.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer errorFile.Close()

	infoHandler := slog.NewTextHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo})
	warnHandler := slog.NewTextHandler(warnFile, &slog.HandlerOptions{Level: slog.LevelWarn})
	errorHandler := slog.NewTextHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError})

	handler := NewLevelHandler(map[slog.Level]slog.Handler{
		slog.LevelInfo:  infoHandler,
		slog.LevelWarn:  warnHandler,
		slog.LevelError: errorHandler,
	}).WithAttrs([]slog.Attr{
		slog.Int64("gid", gid.ID()),
	})

	logger := slog.New(handler).WithGroup("program_info")

	logger.Info("info text")
	logger.Warn("warn text")
	logger.Error("error text")
}
