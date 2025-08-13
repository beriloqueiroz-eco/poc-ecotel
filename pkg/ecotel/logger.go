package ecotel

import (
	"context"

	"github.com/tradersclub/encelado-utilities-go/logger"
	"go.opentelemetry.io/otel/trace"
)

var logServiceName string

func SetLogServiceName(name string) {
	logServiceName = name
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return
	}
	logger.Info(msg,
		append(fields,
			"traceId", spanCtx.TraceID().String(),
			"spanID", spanCtx.SpanID().String(),
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
		)...,
	)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return
	}
	logger.Error(msg,
		append(fields,
			"traceId", spanCtx.TraceID().String(),
			"spanID", spanCtx.SpanID().String(),
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
		)...,
	)
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return
	}
	logger.Debug(msg,
		append(fields,
			"traceId", spanCtx.TraceID().String(),
			"spanID", spanCtx.SpanID().String(),
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
		)...,
	)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return
	}
	logger.Warn(msg,
		append(fields,
			"traceId", spanCtx.TraceID().String(),
			"spanID", spanCtx.SpanID().String(),
			"service.name", logServiceName,
			"resource.service.name", logServiceName,
		)...,
	)
}
