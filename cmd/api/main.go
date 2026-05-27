package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luken/notes-devops-playground/internal/config"
	"github.com/luken/notes-devops-playground/internal/handlers"
	"github.com/luken/notes-devops-playground/internal/logger"
	appMiddleware "github.com/luken/notes-devops-playground/internal/middleware"
	"github.com/luken/notes-devops-playground/internal/repository"
	"github.com/luken/notes-devops-playground/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	log.Info("connected to database")

	repo := repository.NewNoteRepository(pool)
	svc := service.NewNoteService(repo, &service.DefaultValidator{})
	noteHandler := handlers.NewNoteHandler(svc)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(appMiddleware.Recovery(log))
	r.Use(appMiddleware.Logging(log))
	r.Use(appMiddleware.Metrics)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	r.Get("/health", handlers.HealthHandler())
	r.Handle("/metrics", promhttp.Handler())

	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", noteHandler.GetAll)
		r.Post("/", noteHandler.Create)
		r.Get("/{id}", noteHandler.GetByID)
		r.Put("/{id}", noteHandler.Update)
		r.Delete("/{id}", noteHandler.Delete)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server exited gracefully")
}
