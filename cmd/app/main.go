package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MrGreenboom/go-task-manager/internal/handler"
	"github.com/MrGreenboom/go-task-manager/internal/repository"
	"github.com/MrGreenboom/go-task-manager/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is required")
	}

	// JWT_SECRET обязателен для auth middleware/login
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// repositories
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	// services
	taskSvc := service.NewTaskService(taskRepo)
	authSvc := service.NewAuthService(userRepo)

	// handlers
	taskHandler := handler.NewTaskHandler(taskSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// auth routes (public)
	authHandler.RegisterRoutes(mux)

	// tasks routes (protected by JWT middleware)
	protectedMux := http.NewServeMux()
	taskHandler.RegisterRoutes(protectedMux)

	mux.Handle("/tasks", handler.AuthMiddleware(protectedMux))
	mux.Handle("/tasks/", handler.AuthMiddleware(protectedMux))

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
