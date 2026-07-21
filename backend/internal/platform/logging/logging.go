package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, h := range m.handlers {
		if !h.Enabled(ctx, record.Level) {
			continue
		}
		if err := h.Handle(ctx, record.Clone()); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		next[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: next}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	next := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		next[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: next}
}

type levelFilterHandler struct {
	minLevel slog.Level
	maxLevel slog.Level
	inner    slog.Handler
}

func (h *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < h.minLevel || level > h.maxLevel {
		return false
	}
	return h.inner.Enabled(ctx, level)
}

func (h *levelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level < h.minLevel || record.Level > h.maxLevel {
		return nil
	}
	return h.inner.Handle(ctx, record)
}

func (h *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{
		minLevel: h.minLevel,
		maxLevel: h.maxLevel,
		inner:    h.inner.WithAttrs(attrs),
	}
}

func (h *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{
		minLevel: h.minLevel,
		maxLevel: h.maxLevel,
		inner:    h.inner.WithGroup(name),
	}
}

type fileSink struct {
	mu    sync.Mutex
	files []*os.File
}

func (s *fileSink) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, f := range s.files {
		_ = f.Close()
	}
	s.files = nil
}

func New(env string) (*slog.Logger, func(), error) {
	dir, err := resolveLogDir()
	if err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create log dir: %w", err)
	}

	sink := &fileSink{}
	debugFile, err := openLogFile(dir, "debug.log")
	if err != nil {
		return nil, nil, err
	}
	infoFile, err := openLogFile(dir, "info.log")
	if err != nil {
		_ = debugFile.Close()
		return nil, nil, err
	}
	errorFile, err := openLogFile(dir, "error.log")
	if err != nil {
		_ = debugFile.Close()
		_ = infoFile.Close()
		return nil, nil, err
	}
	sink.files = []*os.File{debugFile, infoFile, errorFile}

	consoleLevel := slog.LevelInfo
	if env == "development" {
		consoleLevel = slog.LevelDebug
	}

	handlers := []slog.Handler{
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: consoleLevel}),
		&levelFilterHandler{
			minLevel: slog.LevelDebug,
			maxLevel: slog.LevelDebug,
			inner:    slog.NewJSONHandler(debugFile, &slog.HandlerOptions{Level: slog.LevelDebug}),
		},
		&levelFilterHandler{
			minLevel: slog.LevelInfo,
			maxLevel: slog.LevelWarn,
			inner:    slog.NewJSONHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo}),
		},
		&levelFilterHandler{
			minLevel: slog.LevelError,
			maxLevel: slog.LevelError + 4,
			inner:    slog.NewJSONHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError}),
		},
	}

	logger := slog.New(&multiHandler{handlers: handlers})
	return logger, sink.Close, nil
}

func openLogFile(dir, name string) (*os.File, error) {
	path := filepath.Join(dir, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return f, nil
}

func resolveLogDir() (string, error) {
	if v := os.Getenv("LOG_DIR"); v != "" {
		return v, nil
	}
	for _, candidate := range []string{"logs", "../logs"} {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		parent := filepath.Dir(abs)
		if _, err := os.Stat(filepath.Join(parent, "backend")); err == nil {
			return abs, nil
		}
		if _, err := os.Stat(filepath.Join(parent, "config")); err == nil {
			return abs, nil
		}
	}
	return filepath.Abs("logs")
}
