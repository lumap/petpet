package lib

import "log/slog"

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

func LogError(msg string, args ...any) {
	slog.Error(msg, args...)
}
