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
