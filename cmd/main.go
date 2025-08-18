package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	log "github.com/tradersclub/encelado-utilities-go/logger"
	config "github.com/tradersclub/poc-ecotel/configs"
	"github.com/tradersclub/poc-ecotel/internal"
	"github.com/tradersclub/poc-ecotel/pkg/ecotel"
)

func main() {
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	serviceName := configs.ServiceName
	collectorUrl := configs.OtelExporterEndpoint
	portGin := configs.WebServerPort
	insecure := configs.InsecureOtelCollector
	logPath := configs.LogPath
	logFormat := configs.LogFormat

	// if logFormat == "json" {
	// 	log.DefaultEncoding = "json"
	// }
	logger := log.NewLogger(false, "info")
	ecotel.SetLogServiceName(serviceName)
	ecotel.SetLogger(logger)
	if logFormat == "json" {
		ecotel.UseSlog()
	}
	ctx := context.Background()

	err = redirectStdoutStderr(logPath)
	if err != nil {
		ecotel.Error(ctx, "Error redirecting stdout/stderr:", err)
	}

	// ----------trace--------------
	otelTracer := ecotel.NewTraceEcotel(collectorUrl, serviceName)
	otelTracer.InitTracerProvider(ctx, insecure)
	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	ginServer.Use(otelTracer.GinMiddleware())
	http.DefaultClient = otelTracer.Client

	// ----------metrics--------------
	otelMetrics := ecotel.NewMetricEcotel(collectorUrl, serviceName)
	otelMetrics.InitMeterProvider(insecure)
	ginServer.Use(otelMetrics.GinMetricsMiddleware())

	url := "postgres://user:password@db:5432/database?sslmode=disable"

	pool, err := otelTracer.InitOtelPgxTracer(ctx, url)
	if err != nil {
		ecotel.Error(ctx, "Error initializing database connection pool:", err)
		panic(err)
	}
	db, err := internal.NewDB(ctx, pool)
	if err != nil {
		ecotel.Error(ctx, "Error connecting to database:", err)
		panic(err)
	}
	helloHandler := &internal.HelloHandler{
		IsEnd:        configs.IsEnd,
		Delay:        configs.TestDelay,
		ServiceName:  configs.ServiceName,
		ServiceUrlTo: configs.ServiceBUrl,
		Repo:         internal.NewRepository(db),
	}

	ginServer.GET("/hello", helloHandler.Handle)

	func() {
		ecotel.Info(ctx, fmt.Sprintf("Starting Gin server on port %s", portGin))
		if err := ginServer.Run(portGin); err != nil {
			ecotel.Error(ctx, "Failed to start Gin server:", err)
		}
	}()
}

func redirectStdoutStderr(logFile string) error {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	// Redireciona stdout (fd 1) e stderr (fd 2)
	syscall.Dup2(int(f.Fd()), 1)
	syscall.Dup2(int(f.Fd()), 2)
	return nil
}
