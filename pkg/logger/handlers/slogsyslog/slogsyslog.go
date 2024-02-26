package slogsyslog

import (
	"context"
	"log/slog"
	"log/syslog"
)

type syslogHandler struct {
	h slog.Handler
}

// Custom slog.Handler to write all records to linux syslog daemon
func NewSyslogHandler(prefix string, opts *slog.HandlerOptions) (*syslogHandler, error) {
	s, err := syslog.New(syslog.LOG_DEBUG, prefix)
	if err != nil {
		return &syslogHandler{}, err
	}

	return &syslogHandler{
		h: slog.NewTextHandler(s, opts),
	}, nil
}

func (s *syslogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.h.Enabled(ctx, level)
}
func (s *syslogHandler) Handle(ctx context.Context, record slog.Record) error {
	return s.h.Handle(ctx, record)
}
func (s *syslogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return s.h.WithAttrs(attrs)
}
func (s *syslogHandler) WithGroup(name string) slog.Handler {
	return s.h.WithGroup(name)
}
