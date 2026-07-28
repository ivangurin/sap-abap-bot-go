package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/trace"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

const (
	FieldTraceID   = "trace_id"
	FieldSpanID    = "span_id"
	FieldRequestID = "request_id"
	FieldOrgID     = "org_id"
	FieldUserID    = "user_id"
	FieldIP        = "ip"
)

type Logger interface {
	With(key, value string) Logger

	Debug(ctx context.Context, msg string, values ...any)
	Debugf(ctx context.Context, msg string, args ...any)

	Info(ctx context.Context, msg string, values ...any)
	Infof(ctx context.Context, msg string, args ...any)

	Warn(ctx context.Context, msg string, values ...any)
	Warnf(ctx context.Context, msg string, args ...any)

	Error(ctx context.Context, msg string, values ...any)
	Errorf(ctx context.Context, msg string, args ...any)

	IsWithDebug() bool
}

type logger struct {
	logLevel string
	logger   *slog.Logger
}

func NewLogger(logLevel, logFile string) Logger {
	log := setupPrettySlog(logLevel, logFile)
	slog.SetDefault(log)

	return &logger{
		logLevel: logLevel,
		logger:   log,
	}
}

func setupPrettySlog(logLevel, logFile string) *slog.Logger {
	var output io.Writer
	slogLevel := levelToSlogLevel(logLevel)
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slogLevel,
		},
	}
	if logFile != "" {
		f, err := os.OpenFile(filepath.Clean(logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666) //nolint:gosec
		if err != nil {
			panic(err)
		}
		output = f
	} else {
		output = os.Stdout
	}

	handler := opts.NewPrettyHandler(output)

	return slog.New(handler)
}

func levelToSlogLevel(logLevel string) slog.Level {
	switch logLevel {
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func (l *logger) With(key, value string) Logger {
	return &logger{
		logLevel: l.logLevel,
		logger:   l.logger.With(slog.String(key, value)),
	}
}

func (l *logger) Debug(ctx context.Context, msg string, values ...any) {
	l.logger.Debug(msg, l.getFields(ctx, values)...)
}

func (l *logger) Debugf(ctx context.Context, msg string, args ...any) {
	l.logger.Debug(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) Info(ctx context.Context, msg string, values ...any) {
	l.logger.Info(msg, l.getFields(ctx, values)...)
}

func (l *logger) Infof(ctx context.Context, msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) Warn(ctx context.Context, msg string, values ...any) {
	l.logger.Warn(msg, l.getFields(ctx, values)...)
}

func (l *logger) Warnf(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) Error(ctx context.Context, msg string, values ...any) {
	l.logger.Error(msg, l.getFields(ctx, values)...)
}

func (l *logger) Errorf(ctx context.Context, msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) getFields(ctx context.Context, args []any) []any {
	traceID := l.getTraceID(ctx)
	if traceID != "" {
		args = append(args, FieldTraceID, traceID)
	}

	spanID := l.getSpanID(ctx)
	if spanID != "" {
		args = append(args, FieldSpanID, spanID)
	}

	return args
}

func (l *logger) getTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}

func (l *logger) getSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

func (l *logger) IsWithDebug() bool {
	return l.logLevel == LevelDebug
}
