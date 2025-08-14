package ecotel

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"
)

type testLogger struct {
	lastMsg    string
	lastLevel  string
	lastFields []any
}

func (l *testLogger) Info(msg string, fields ...any) {
	l.lastLevel = "info"
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *testLogger) Error(msg string, fields ...any) {
	l.lastLevel = "error"
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *testLogger) Debug(msg string, fields ...any) {
	l.lastLevel = "debug"
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *testLogger) Warn(msg string, fields ...any) {
	l.lastLevel = "warn"
	l.lastMsg = msg
	l.lastFields = fields
}
func (l *testLogger) Fatal(msg string, fields ...any) {
	l.lastLevel = "fatal"
	l.lastMsg = msg
	l.lastFields = fields
}

func TestLoggerWithContextWithTraceOtelProvider(t *testing.T) {
	l := &testLogger{}
	SetLogger(l)
	SetLogServiceName("service-test")

	tracer := trace.NewTracerProvider().Tracer("test-tracer")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	Info(ctx, "info-msg", "foo", "bar")
	if l.lastLevel != "info" || l.lastMsg != "info-msg" {
		t.Errorf("Info log not called correctly")
	}
	foundTrace := false
	foundService := false
	for i := 0; i < len(l.lastFields); i += 2 {
		if l.lastFields[i] == "traceId" && l.lastFields[i+1] != "" {
			foundTrace = true
		}
		if l.lastFields[i] == "service.name" && l.lastFields[i+1] == "service-test" {
			foundService = true
		}
	}
	if !foundTrace || !foundService {
		t.Errorf("traceId or service.name not found in fields")
	}

	Error(ctx, "error-msg")
	if l.lastLevel != "error" || l.lastMsg != "error-msg" {
		t.Errorf("Error log not called correctly")
	}

	Debug(ctx, "debug-msg")
	if l.lastLevel != "debug" || l.lastMsg != "debug-msg" {
		t.Errorf("Debug log not called correctly")
	}

	Warn(ctx, "warn-msg")
	if l.lastLevel != "warn" || l.lastMsg != "warn-msg" {
		t.Errorf("Warn log not called correctly")
	}
}

func TestLoggerWithContextWithParametersByContextValue(t *testing.T) {
	l := &testLogger{}
	SetLogger(l)
	SetLogServiceName("service-test")
	ctx := context.WithValue(context.Background(), "traceId", "blablabla")

	Info(ctx, "info-msg", "foo", "bar")
	if l.lastLevel != "info" || l.lastMsg != "info-msg" {
		t.Errorf("Info log not called correctly")
	}
	foundTrace := false
	foundService := false
	for i := 0; i < len(l.lastFields); i += 2 {
		if l.lastFields[i] == "traceId" && l.lastFields[i+1] != "" {
			foundTrace = true
		}
		if l.lastFields[i] == "service.name" && l.lastFields[i+1] == "service-test" {
			foundService = true
		}
	}
	if !foundTrace || !foundService {
		t.Errorf("traceId or service.name not found in fields")
	}

	Error(ctx, "error-msg")
	if l.lastLevel != "error" || l.lastMsg != "error-msg" {
		t.Errorf("Error log not called correctly")
	}

	Debug(ctx, "debug-msg")
	if l.lastLevel != "debug" || l.lastMsg != "debug-msg" {
		t.Errorf("Debug log not called correctly")
	}

	Warn(ctx, "warn-msg")
	if l.lastLevel != "warn" || l.lastMsg != "warn-msg" {
		t.Errorf("Warn log not called correctly")
	}
}
