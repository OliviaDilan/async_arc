package handler

import (
	"context"
	"errors"
	"log"

	"github.com/OliviaDilan/async_arc/accounting/internal/account"
	"github.com/OliviaDilan/async_arc/pkg/auth"
	contractAuth "github.com/OliviaDilan/async_arc/pkg/contract/auth"
)

type Handler struct {
	accountRepo account.Repository
	authClient auth.Client
}

func NewHandler(accountRepo account.Repository, authClient auth.Client) *Handler {
	return &Handler{
		accountRepo: accountRepo,
		authClient: authClient,
	}
}

func (h *Handler) OnUserCreated(ctx context.Context, body []byte) error {
	var message contractAuth.UserCreatedV1

	err := message.Unmarshal(body)
	if err != nil {
		return err
	}

	if message.Username == "" {
		return errors.New("invalid username")
	}

	err = h.accountRepo.Create(message.Username); 
	if err != nil {
		return err
	}

	log.Printf("Account %s created", message.Username)

	return nil
}