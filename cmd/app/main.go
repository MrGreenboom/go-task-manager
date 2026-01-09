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

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	taskHandler.RegisterRoutes(mux)

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
