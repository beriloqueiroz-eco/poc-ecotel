package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tradersclub/encelado-utilities-go/logger"
	config "github.com/tradersclub/poc-ecotel/configs"
	"github.com/tradersclub/poc-ecotel/internal"
	"github.com/tradersclub/poc-ecotel/pkg/ecotel"
)

type ctxKey string

func main() {
	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}

	serviceName := configs.ServiceName
	collectorUrl := configs.OtelExporterEndpoint
	portGin := configs.WebServerPort
	insecure := configs.InsecureOtelCollector

	log := logger.NewZapLogger(true, "info")
	logger.SetGlobal(log)

	ecotel.SetLogServiceName(serviceName)
	ecotel.SetLogger(log)

	ctx := context.Background()

	err = redirectStdoutStderr("/var/log/app.log")
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

	helloHandler := &internal.HelloHandler{
		IsEnd:        configs.IsEnd,
		Delay:        configs.TestDelay,
		ServiceName:  configs.ServiceName,
		ServiceUrlTo: configs.ServiceBUrl,
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
