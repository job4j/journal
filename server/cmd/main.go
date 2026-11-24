package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi "journal/server/internal/api"
	"journal/server/internal/config"
	"journal/server/internal/repository"
	"journal/server/internal/service"
	"journal/server/internal/storage"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	txManager := repository.NewTransactionManager(pool)
	authRepo := repository.NewAuthRepository()
	authService := service.NewAuthService(txManager, authRepo, 24*time.Hour)
	appRepository := repository.NewRepository()
	roleService := service.NewRoleService(txManager, appRepository)
	userService := service.NewUserService(txManager, appRepository)
	academicYearService := service.NewAcademicYearService(txManager, appRepository)
	subjectService := service.NewSubjectService(txManager, appRepository)
	classService := service.NewClassService(txManager, appRepository)
	var objectStore service.AttachmentObjectStore
	if cfg.Storage.Endpoint != "" {
		store, storageErr := storage.NewMinio(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey, cfg.Storage.Bucket, cfg.Storage.Secure)
		if storageErr != nil {
			logger.Error("create file storage", "error", storageErr)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		storageErr = store.EnsureBucket(ctx)
		cancel()
		if storageErr != nil {
			logger.Error("initialize file storage bucket", "error", storageErr)
			os.Exit(1)
		}
		objectStore = store
	}
	app := fiber.New(fiber.Config{BodyLimit: 24 << 20, ErrorHandler: func(c *fiber.Ctx, err error) error {
		if errors.Is(err, fiber.ErrRequestEntityTooLarge) {
			return c.Status(413).JSON(fiber.Map{"code": "request_too_large", "message": "Файл слишком большой. Максимум — 20 МБ"})
		}
		logger.Error("http request failed", "request_id", c.GetRespHeader("X-Request-ID"), "method", c.Method(), "path", c.Path(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "internal_error", "message": "Внутренняя ошибка сервера"})
	}})
	app.Use(httpapi.RequestLog(logger))
	apiRouter := app.Group("/api/v1")
	httpapi.NewAuthHandler(authService, logger, cfg.CookieSecure).Register(apiRouter)
	httpapi.NewRoleHandler(roleService, logger).Register(apiRouter)
	httpapi.NewUserHandler(userService, logger).Register(apiRouter)
	httpapi.NewAcademicYearHandler(academicYearService, logger).Register(apiRouter)
	httpapi.NewSubjectHandler(subjectService, logger).Register(apiRouter)
	httpapi.NewClassHandler(classService, logger).Register(apiRouter)
	httpapi.NewAttachmentHandler(service.NewAttachmentService(txManager, appRepository, objectStore), logger).Register(apiRouter)
	httpapi.NewScheduleHandler(service.NewScheduleService(txManager, appRepository), logger).Register(apiRouter)

	// The root spec is served by the Swagger middleware at /api/openapi.yml.
	// Expose its sibling path files so relative $ref values resolve in Swagger UI.
	app.Static("/api/paths", "./api/paths")
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./api/openapi.yml",
		Path:     "docs",
		Title:    "API documentation",
	}))
	listenErrors := make(chan error, 1)
	go func() {
		listenErrors <- app.Listen(cfg.HTTPAddress)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-listenErrors:
		logger.Error("listen HTTP", "error", err)
		return
	case <-ctx.Done():
	}
	logger.Info("shutting down HTTP server", "timeout", cfg.ShutdownTimeout)
	if err := app.ShutdownWithTimeout(cfg.ShutdownTimeout); err != nil {
		logger.Error("shutdown http server", "error", err)
	}
}
