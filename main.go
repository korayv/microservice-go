package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"microservicetest/pkg/config"
	_ "microservicetest/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()

	app := fiber.New()

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Get("/", func(c *fiber.Ctx) error {
		zap.L().Info("App starting!")
		return c.SendString("Hello, World!")
	})
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()
	zap.L().Info("Server started on port ", zap.String("port", appConfig.Port))
	gracefulShutdown(app)
}
func gracefulShutdown(app *fiber.App) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("Shutting down gracefully")
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Failed to shutdown server", zap.Error(err))
	}
	zap.L().Info("Server gracefully stopped")
}
