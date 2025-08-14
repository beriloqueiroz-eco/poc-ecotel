package ecotel

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc/credentials"
)

type TraceEcotel struct {
	Client       *http.Client
	collectorUrl string
	serviceName  string
}

func NewTraceEcotel(collectorUrl, serviceName string) *TraceEcotel {
	return &TraceEcotel{
		collectorUrl: collectorUrl,
		serviceName:  serviceName,
		Client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (o *TraceEcotel) InitTracerProvider(ctx context.Context, insecure bool) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}
	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(o.collectorUrl),
		),
	)
	if err != nil {
		return nil, err
	}
	// Create the resource to be traced
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(o.serviceName),
		semconv.ServiceVersion("v0.0.1"),
	)

	// Configure the trace provider
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(2*time.Second)),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	return traceProvider.Shutdown, nil
}

// Ou, se quiser customizar o transporte:
func (o *TraceEcotel) NewInstrumentedClientWithTransport(transport http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(transport),
	}
}

func (o *TraceEcotel) GinMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(o.serviceName)
}

func (o *TraceEcotel) InitOtelPgxTracer(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := otelpgx.RecordStats(conn); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	return conn, nil

}

func MainExample() {
	serviceName := "service-example"
	collectorUrl := "http://localhost:4317"
	portGin := ":8080"
	insecure := true

	ctx := context.Background()

	otelTracer := NewTraceEcotel(collectorUrl, serviceName)
	shutdown, err := otelTracer.InitTracerProvider(ctx, insecure)
	if err != nil {
		log.Fatalf("Failed to initialize tracer provider: %v", err)
	}
	defer shutdown(ctx)
	ginServer := gin.Default()
	ginServer.Use(otelTracer.GinMiddleware())
	http.DefaultClient = otelTracer.Client

	ginServer.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	func() {
		fmt.Println("Starting Gin server on", portGin)
		if err := ginServer.Run(portGin); err != nil {
			log.Fatalf("Gin server error: %v", err)
		}
	}()
}
