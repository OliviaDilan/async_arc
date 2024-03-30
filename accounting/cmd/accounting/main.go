package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/OliviaDilan/async_arc/pkg/amqp"
	"github.com/OliviaDilan/async_arc/pkg/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/OliviaDilan/async_arc/accounting/internal/account"
	internalAMQP "github.com/OliviaDilan/async_arc/accounting/internal/amqp"
	"github.com/OliviaDilan/async_arc/accounting/internal/config"
	"github.com/OliviaDilan/async_arc/accounting/internal/handler"
	"github.com/OliviaDilan/async_arc/accounting/internal/task"
)

func main() {
	ctx, shutdown := appContext()
	defer shutdown()

	cfg := config.Server{}

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatal(err)
	}

	amqpClient, err := amqp.NewClient(cfg.AMQP.URI())
	if err != nil {
		log.Fatal(err)
	}

	consumerSet := internalAMQP.NewConsumerSet(amqpClient)

	accountRepo := account.NewInMemoryRepository()

	taskRepo := task.NewInMemoryRepository()

	authClient := auth.NewClient(cfg.Auth.Host, cfg.Auth.Port)

	h := handler.NewHandler(accountRepo, taskRepo, authClient)

	err = consumerSet.Subscribe("user_created", h.OnUserCreated)
	if err != nil {
		log.Fatal(err)
	}

	err = consumerSet.Subscribe("task_created", h.OnTaskCreated)
	if err != nil {
		log.Fatal(err)
	}

	err = consumerSet.Subscribe("task_assigned", h.OnTaskAssigned)
	if err != nil {
		log.Fatal(err)
	}

	err = consumerSet.Subscribe("task_completed", h.OnTaskCompleted)
	if err != nil {
		log.Fatal(err)
	}

	consumerSet.StartConsumers(ctx)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(auth.Middleware(authClient))
	r.Get("/all_accounts", h.GetAccounts)


	path := cfg.Host + ":" + cfg.Port

	srv := &http.Server{
		Addr:    path,
		Handler: r,
	}

	done := make(chan error, 1)

	go func() {
		<-ctx.Done()
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

func appContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}