package ecotel

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Fatal(msg string, fields ...any)
}

type slogAdapter struct {
	slog *slog.Logger
}

func (s *slogAdapter) Info(msg string, fields ...any)  { s.slog.Info(msg, fields...) }
func (s *slogAdapter) Error(msg string, fields ...any) { s.slog.Error(msg, fields...) }
func (s *slogAdapter) Debug(msg string, fields ...any) { s.slog.Debug(msg, fields...) }
func (s *slogAdapter) Warn(msg string, fields ...any)  { s.slog.Warn(msg, fields...) }
func (s *slogAdapter) Fatal(msg string, fields ...any) { s.slog.Error(msg, fields...); os.Exit(1) }

var logger Logger

var logServiceName string

func SetLogServiceName(name string) {
	logServiceName = name
}

func SetLogger(l Logger) {
	logger = l
}

func UseSlog() {
	s := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(s)
	logger = &slogAdapter{slog: s}
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	traceIdCtx := ctx.Value("traceId")
	spanIdCtx := ctx.Value("spanId")

	traceId := ""
	spanId := ""
	if spanIdCtx != nil {
		spanId = spanIdCtx.(string)
	}
	if traceIdCtx != nil {
		traceId = traceIdCtx.(string)
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		if spanId == "" {
			spanId = spanCtx.SpanID().String()
		}
		if traceId == "" {
			traceId = spanCtx.TraceID().String()
		}
	}
	logger.Info(msg,
		append(fields,
			"traceId", traceId,
			"spanId", spanId,
			"service_name", logServiceName,
		)...,
	)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	traceIdCtx := ctx.Value("traceId")
	spanIdCtx := ctx.Value("spanId")

	traceId := ""
	spanId := ""
	if spanIdCtx != nil {
		spanId = spanIdCtx.(string)
	}
	if traceIdCtx != nil {
		traceId = traceIdCtx.(string)
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		if spanId == "" {
			spanId = spanCtx.SpanID().String()
		}
		if traceId == "" {
			traceId = spanCtx.TraceID().String()
		}
	}
	logger.Error(msg,
		append(fields,
			"traceId", traceId,
			"spanId", spanId,
			"service_name", logServiceName,
		)...,
	)
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	traceIdCtx := ctx.Value("traceId")
	spanIdCtx := ctx.Value("spanId")

	traceId := ""
	spanId := ""
	if spanIdCtx != nil {
		spanId = spanIdCtx.(string)
	}
	if traceIdCtx != nil {
		traceId = traceIdCtx.(string)
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		if spanId == "" {
			spanId = spanCtx.SpanID().String()
		}
		if traceId == "" {
			traceId = spanCtx.TraceID().String()
		}
	}
	logger.Debug(msg,
		append(fields,
			"traceId", traceId,
			"spanId", spanId,
			"service_name", logServiceName,
		)...,
	)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	traceIdCtx := ctx.Value("traceId")
	spanIdCtx := ctx.Value("spanId")

	traceId := ""
	spanId := ""
	if spanIdCtx != nil {
		spanId = spanIdCtx.(string)
	}
	if traceIdCtx != nil {
		traceId = traceIdCtx.(string)
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		if spanId == "" {
			spanId = spanCtx.SpanID().String()
		}
		if traceId == "" {
			traceId = spanCtx.TraceID().String()
		}
	}
	logger.Warn(msg,
		append(fields,
			"traceId", traceId,
			"spanId", spanId,
			"service_name", logServiceName,
		)...,
	)
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	traceIdCtx := ctx.Value("traceId")
	spanIdCtx := ctx.Value("spanId")

	traceId := ""
	spanId := ""
	if spanIdCtx != nil {
		spanId = spanIdCtx.(string)
	}
	if traceIdCtx != nil {
		traceId = traceIdCtx.(string)
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		if spanId == "" {
			spanId = spanCtx.SpanID().String()
		}
		if traceId == "" {
			traceId = spanCtx.TraceID().String()
		}
	}
	logger.Fatal(msg,
		append(fields,
			"traceId", traceId,
			"spanId", spanId,
			"service_name", logServiceName,
		)...,
	)
}
