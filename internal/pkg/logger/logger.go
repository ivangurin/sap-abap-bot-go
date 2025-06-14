package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Logger struct {
	ctx    context.Context
	logger *zap.Logger
}

const (
	keyTraceID = "trace_id"
	keySpanID  = "span_id"
)

func NewLogger(ctx context.Context, opts ...ConfigOption) *Logger {
	logger, err := NewConfig(opts...).Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		ctx:    ctx,
		logger: logger,
	}
}

func (l *Logger) Close() error {
	err := l.logger.Sync()
	if err != nil {
		return fmt.Errorf("logger sync: %w", err)
	}

	return nil
}

func (l *Logger) Debug(m string) {
	l.logger.Debug(m, getFields(l.ctx)...)
}

func (l *Logger) Debugf(m string, args ...any) {
	l.logger.Debug(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func (l *Logger) Info(m string) {
	l.logger.Info(m, getFields(l.ctx)...)
}

func (l *Logger) Infof(m string, args ...any) {
	l.logger.Info(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func (l *Logger) Warn(m string) {
	l.logger.Warn(m, getFields(l.ctx)...)
}

func (l *Logger) Warnf(m string, args ...any) {
	l.logger.Warn(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func (l *Logger) Error(m string) {
	l.logger.Error(m, getFields(l.ctx)...)
}

func (l *Logger) Errorf(m string, args ...any) {
	l.logger.Error(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func (l *Logger) Panic(m string) {
	l.logger.Panic(m, getFields(l.ctx)...)
}

func (l *Logger) Panicf(m string, args ...any) {
	l.logger.Panic(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func (l *Logger) Fatal(m string) {
	l.logger.Fatal(m, getFields(l.ctx)...)
}

func (l *Logger) Fatalf(m string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(m, args...), getFields(l.ctx)...)
}

func getFields(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 2)
	traceID := getTraceID(ctx)
	if traceID != "" {
		fields = append(fields, zap.String(keyTraceID, traceID))
	}
	spanID := getSpanID(ctx)
	if spanID != "" {
		fields = append(fields, zap.String(keySpanID, spanID))
	}
	return fields
}

func getTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func getSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}
