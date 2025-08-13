package ecotel

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

type MetricEcotel struct {
	collectorUrl   string
	serviceName    string
	metricProvider *sdkmetric.MeterProvider
}

func NewMetricEcotel(collectorUrl, serviceName string) *MetricEcotel {
	return &MetricEcotel{
		collectorUrl: collectorUrl,
		serviceName:  serviceName,
	}
}

func (o *MetricEcotel) InitMeterProvider(insecure bool) error {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlpmetricgrpc.WithInsecure()
	}
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithEndpoint(o.collectorUrl),
		secureOption,
	)
	if err != nil {
		return err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter),
		),
		sdkmetric.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(o.serviceName),
			),
		),
	)
	otel.SetMeterProvider(mp)
	o.metricProvider = mp
	return nil
}

func (o *MetricEcotel) GinMetricsMiddleware() gin.HandlerFunc {
	meter := o.metricProvider.Meter(o.serviceName)
	requestCounter, _ := meter.Int64Counter("http_requests_total")
	latencyHistogram, _ := meter.Float64Histogram("http_request_duration_seconds")
	errorCounter, _ := meter.Int64Counter("http_requests_errors_total")
	methodCounter, _ := meter.Int64Counter("http_requests_by_method_total")
	routeCounter, _ := meter.Int64Counter("http_requests_by_route_total")

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		// Extrai trace_id e span_id do contexto, se houver
		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		traceID := ""
		spanID := ""
		if spanCtx.HasTraceID() {
			traceID = spanCtx.TraceID().String()
		}
		if spanCtx.HasSpanID() {
			spanID = spanCtx.SpanID().String()
		}

		attrs := []attribute.KeyValue{
			attribute.String("method", c.Request.Method),
			attribute.String("route", c.FullPath()),
			attribute.Int("status", c.Writer.Status()),
		}
		if traceID != "" {
			attrs = append(attrs, attribute.String("trace_id", traceID))
		}
		if spanID != "" {
			attrs = append(attrs, attribute.String("span_id", spanID))
		}

		requestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		latencyHistogram.Record(c.Request.Context(), duration, metric.WithAttributes(attrs...))

		methodAttrs := []attribute.KeyValue{attribute.String("method", c.Request.Method)}
		if traceID != "" {
			methodAttrs = append(methodAttrs, attribute.String("trace_id", traceID))
		}
		if spanID != "" {
			methodAttrs = append(methodAttrs, attribute.String("span_id", spanID))
		}
		methodCounter.Add(c.Request.Context(), 1, metric.WithAttributes(methodAttrs...))

		routeAttrs := []attribute.KeyValue{attribute.String("route", c.FullPath())}
		if traceID != "" {
			routeAttrs = append(routeAttrs, attribute.String("trace_id", traceID))
		}
		if spanID != "" {
			routeAttrs = append(routeAttrs, attribute.String("span_id", spanID))
		}
		routeCounter.Add(c.Request.Context(), 1, metric.WithAttributes(routeAttrs...))

		if c.Writer.Status() >= 400 {
			errorCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		}
	}
}
