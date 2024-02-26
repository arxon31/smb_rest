package logger

import (
	"git.spbec-mining.ru/arxon31/sambaMW/pkg/logger/handlers/slogsyslog"
	"log/slog"
	"os"
	"strings"
)

const (
	envDev   = "dev"
	envStage = "stage"
	envProd  = "prod"
)

func New(appName, env string) *slog.Logger {
	level := levelFromEnv(env)

	switch strings.ToLower(env) {
	case envDev:
		return newLogger(appName, level)
	case envStage:
		return newSyslogLogger(appName, level)
	case envProd:
		return newSyslogLogger(appName, level)
	default:
		return newLogger(appName, level)

	}
}

func newSyslogLogger(appName string, level slog.Level) *slog.Logger {

	syslogHandler, err := slogsyslog.NewSyslogHandler(
		appName,
		&slog.HandlerOptions{Level: level})
	if err != nil {
		panic(err)
	}

	return slog.New(syslogHandler)

}

func newLogger(appName string, level slog.Level) *slog.Logger {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	logger := slog.New(textHandler)

	return logger
}
func levelFromEnv(env string) slog.Level {
	var l slog.Level

	switch strings.ToLower(env) {
	case envDev:
		l = slog.LevelDebug
	case envStage:
		l = slog.LevelInfo
	case envProd:
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	return l
}
