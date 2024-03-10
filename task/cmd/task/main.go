package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OliviaDilan/async_arc/pkg/auth"
	"github.com/OliviaDilan/async_arc/pkg/amqp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/OliviaDilan/async_arc/task/internal/config"
	"github.com/OliviaDilan/async_arc/task/internal/handler"
	"github.com/OliviaDilan/async_arc/task/internal/task"
)

func main() {
	// Handle HTTP request
	cfg := config.Server{}

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatal(err)
	}

	amqpClient, err := amqp.NewClient(cfg.AMQP.URI())
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := task.NewInMemoryRepository()

	authClient := auth.NewClient(cfg.Auth.Host, cfg.Auth.Port)

	h := handler.NewHandler(taskRepo, authClient)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.Middleware(authClient))

	r.Post("/create_task", h.CreateTask)
	r.Post("/assign_all_tasks", h.AssignTasks)
	r.Post("/close_task", h.CloseTask)
	r.Get("/all_tasks", h.GetTasks)
	r.Get("/all_tasks_by_assignee", h.GetTasksByAssignee)

	path := cfg.Host + ":" + cfg.Port

	srv := &http.Server{
		Addr:    path,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)

	go func() {
		<-quit
		log.Print("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := srv.Shutdown(ctx)
		amqpClient.Close()
		done <- err
	}()

	log.Printf("Listening on %s", path)
	_ = srv.ListenAndServe()

	err = <-done
	log.Print("Stopping server with error: ", err)
}