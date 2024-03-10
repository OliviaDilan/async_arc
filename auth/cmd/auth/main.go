package main

import (
	"log"
	"net/http"
	"time"
	"context"
	"os"

	"github.com/go-chi/chi"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/OliviaDilan/async_arc/pkg/amqp"

	internalAMQP "github.com/OliviaDilan/async_arc/auth/internal/amqp"

	"github.com/OliviaDilan/async_arc/auth/internal/config"
	"github.com/OliviaDilan/async_arc/auth/internal/handler"
	"github.com/OliviaDilan/async_arc/auth/internal/user"
	"github.com/OliviaDilan/async_arc/auth/internal/jwt"
)

// main.go
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

	publisherSet, err := internalAMQP.NewPublisherSet(amqpClient)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewInMemoryRepository()
	jwtService := jwt.NewService(cfg.JWT.Secret)

	h := handler.NewHandler(userRepo, jwtService, publisherSet)

	r := chi.NewRouter()

	r.Get("/users", h.GetUsers)
	r.Post("/registration", h.Registration)
	r.Post("/login", h.Login)
	r.Post("/decode_token", h.DecodeToken)

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
