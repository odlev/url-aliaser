package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odlev/url-aliaser/internal/config"
	"github.com/odlev/url-aliaser/internal/handlers"
	auth "github.com/odlev/url-aliaser/internal/lib/auth"
	"github.com/odlev/url-aliaser/internal/lib/sl"
	"github.com/odlev/url-aliaser/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("startings url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	s, err := storage.NewSQLite(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := gin.Default()

	//router.Use(gin.Logger())
	router.POST("/login", auth.LoginHandler(*cfg, log))
	authGroup := router.Group("/auth")
	authGroup.Use(auth.AuthMiddleware(log))
	{
		authGroup.POST("/save", handlers.NewSave(log, s))
		authGroup.DELETE("/delete", handlers.NewDelete(log, s))
	}

	router.GET("/:alias", handlers.NewRedirect(log, s))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	}()
	
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Info("Server forsed to shutdownd:", sl.Err(err))
	}
	log.Info("server exiting")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
