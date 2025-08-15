package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

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
	logPath := configs.LogPath

	ecotel.SetLogServiceName(serviceName)
	ecotel.UseSlog()

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
	dbConfig := internal.DBConfig{
		Host:     "db",
		Port:     5432,
		Database: "database",
		User:     "user",
		Password: "password",
		SSLMode:  "disable",
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.SSLMode)

	pool, err := otelTracer.InitOtelPgxTracer(ctx, url)
	if err != nil {
		ecotel.Error(ctx, "Error initializing database connection pool:", err)
		panic(err)
	}
	db, err := internal.NewDB(ctx, dbConfig, pool)
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
