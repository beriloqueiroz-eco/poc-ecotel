package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	ctx := context.Background()

	// ----------trace--------------
	otelTracer := ecotel.NewTraceEcotel(collectorUrl, serviceName)
	otelTracer.InitTracerProvider(ctx, insecure)
	ginServer := gin.Default()
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
		fmt.Println("Starting Gin server on", portGin)
		if err := ginServer.Run(portGin); err != nil {
			log.Fatalf("Gin server error: %v", err)
		}
	}()
}
