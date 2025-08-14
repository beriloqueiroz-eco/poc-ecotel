package ecotel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Fatal(msg string, fields ...any)
}

var logger Logger

var logServiceName string

func SetLogServiceName(name string) {
	logServiceName = name
}

func SetLogger(l Logger) {
	logger = l
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
			"spanID", spanId,
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
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
			"spanID", spanId,
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
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
			"spanID", spanId,
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
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
			"spanID", spanId,
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
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
			"spanID", spanId,
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
		)...,
	)
}
