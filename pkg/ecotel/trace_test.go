package ecotel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func TestInitTracerProvider(t *testing.T) {
	tracer := NewTraceEcotel("localhost:4317", "test-service")
	shutdown, err := tracer.InitTracerProvider(context.Background(), true)
	if err != nil {
		t.Fatalf("failed to init tracer provider: %v", err)
	}
	if shutdown == nil {
		t.Error("shutdown function is nil")
	}
}

func TestGinMiddlewareSpan(t *testing.T) {
	tracer := NewTraceEcotel("localhost:4317", "test-service")
	_, err := tracer.InitTracerProvider(context.Background(), true)
	if err != nil {
		t.Fatalf("failed to init tracer provider: %v", err)
	}
	r := gin.New()
	r.Use(tracer.GinMiddleware())
	r.GET("/test", func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if !span.SpanContext().IsValid() {
			t.Error("span not found in context")
		}
		c.Status(200)
	})

	w := performTraceRequest(r, "GET", "/test")
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// Helper para simular requisições HTTP no Gin
func performTraceRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestGinMiddleware_UsesIncomingTraceparent(t *testing.T) {
	tracer := NewTraceEcotel("localhost:4317", "test-service")
	_, err := tracer.InitTracerProvider(context.Background(), true)
	if err != nil {
		t.Fatalf("failed to init tracer provider: %v", err)
	}
	r := gin.New()
	r.Use(tracer.GinMiddleware())

	var capturedTraceID string
	r.GET("/traceparent", func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		capturedTraceID = span.SpanContext().TraceID().String()
		c.Status(200)
	})

	// Simula um traceparent header vindo do Istio
	traceparent := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/traceparent", nil)
	req.Header.Set("traceparent", traceparent)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	// O traceID do span deve ser igual ao do header (4bf92f3577b34da6a3ce929d0e0e4736)
	if capturedTraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Errorf("traceID not propagated from traceparent header, got %s", capturedTraceID)
	}
}
