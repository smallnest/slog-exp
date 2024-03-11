package slogexp

import (
	"context"
	"errors"
	"log/slog"
)

var ErrBufferFull = errors.New("buffer full")

type bufferedHandler struct {
	slog.Handler
	buffer chan *slog.Record
	block  bool
}

func (h *bufferedHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.block {
		select {
		case h.buffer <- &record:
		}
	} else {
		select {
		case h.buffer <- &record:
		default: //discard the record if the buffer is full
			return ErrBufferFull
		}
	}

	return nil
}

func (h *bufferedHandler) start() error {
	for {
		select {
		case record := <-h.buffer:
			h.Handler.Handle(context.Background(), *record)
		}
	}
}

// Buffer returns a new handler that buffers records and sends them to the
// given handler.
// If block is true, the buffer will block when it is full.
// If block is false, the buffer will discard the record if it is full.
func Buffer(h slog.Handler, bufferSize int, block bool) slog.Handler {
	buffered := &bufferedHandler{
		Handler: h,
		buffer:  make(chan *slog.Record, bufferSize),
		block:   block,
	}

	go buffered.start()

	return buffered
}
