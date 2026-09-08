package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "journal/server/internal/api"
	"journal/server/internal/repository/postgres"
	"journal/server/internal/service"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@127.0.0.1:5433/journal?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Error("create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	authRepo := postgres.NewAuthRepository()
	authService := service.NewAuthService(postgres.NewTransactionManager(pool), authRepo, 24*time.Hour)
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		logger.Error("http request failed", "method", c.Method(), "path", c.Path(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "internal_error", "message": "Внутренняя ошибка сервера"})
	}})
	httpapi.NewAuthHandler(authService, logger, os.Getenv("COOKIE_SECURE") != "false").Register(app.Group("/api/v1"))

	// The root spec is served by the Swagger middleware at /api/openapi.yml.
	// Expose its sibling path files so relative $ref values resolve in Swagger UI.
	app.Static("/api/paths", "./api/paths")
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./api/openapi.yml",
		Path:     "docs",
		Title:    "API documentation",
	}))
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}
	go func() {
		if err := app.Listen(address); err != nil {
			logger.Info("http server stopped", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("shutdown http server", "error", err)
	}
}
