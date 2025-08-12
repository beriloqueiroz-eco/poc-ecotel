package ecotel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInitMeterProvider(t *testing.T) {
	metrics := NewMetricEcotel("localhost:4317", "test-service")
	err := metrics.InitMeterProvider(true)
	if err != nil {
		t.Fatalf("failed to init meter provider: %v", err)
	}
}

func TestGinMetricsMiddleware(t *testing.T) {
	metrics := NewMetricEcotel("localhost:4317", "test-service")
	_ = metrics.InitMeterProvider(true)
	r := gin.New()
	r.Use(metrics.GinMetricsMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	w := performRequest(r, "GET", "/test")
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// Helper para simular requisições HTTP no Gin
func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}
