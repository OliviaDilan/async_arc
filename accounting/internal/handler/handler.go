package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/OliviaDilan/async_arc/accounting/internal/account"
	"github.com/OliviaDilan/async_arc/accounting/internal/task"
	"github.com/OliviaDilan/async_arc/pkg/auth"
	contractAuth "github.com/OliviaDilan/async_arc/pkg/contract/auth"
	contractTask "github.com/OliviaDilan/async_arc/pkg/contract/task"

	"math/rand"
)

type Handler struct {
	accountRepo account.Repository
	taskRepo task.Repository
	authClient auth.Client
}

func NewHandler(accountRepo account.Repository, taskRepo task.Repository, authClient auth.Client) *Handler {
	return &Handler{
		accountRepo: accountRepo,
		taskRepo: taskRepo,
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

func (h *Handler) OnTaskCreated(ctx context.Context, body []byte) error {
	var message contractTask.TaskCreatedV1

	err := message.Unmarshal(body)
	if err != nil {
		return err
	}

	if message.TaskID == 0 {
		return errors.New("invalid taskID")
	}

	taskCost := 10 + rand.Intn(10) // 10 - 20
	taskReward := 20 + rand.Intn(20) // 20 - 40

	estimatedTask, err := h.taskRepo.Create(message.TaskID, taskCost, taskReward);
	if err != nil {
		return err
	}

	log.Printf("Task %d created with cost %d and reward %d", estimatedTask.ID, estimatedTask.Cost, estimatedTask.Reward)

	return nil
}

func (h *Handler) OnTaskAssigned(ctx context.Context, body []byte) error {
	var message contractTask.TaskAssignedV1

	err := message.Unmarshal(body)
	if err != nil {
		return err
	}

	if message.TaskID == 0 {
		return errors.New("invalid taskID")
	}

	if message.UserID == "" {
		return errors.New("invalid assignee")
	}

	task, err := h.taskRepo.GetTaskByExternalID(message.TaskID)
	if err != nil {
		return err
	}

	h.accountRepo.Withdraw(message.UserID, task.Cost);

	log.Printf("Task %d assigned to %s", message.TaskID, message.UserID)

	return nil
}

func (h *Handler) OnTaskCompleted(ctx context.Context, body []byte) error {
	var message contractTask.TaskCompletedV1

	err := message.Unmarshal(body)
	if err != nil {
		return err
	}

	if message.TaskID == 0 {
		return errors.New("invalid taskID")
	}

	if message.UserID == "" {
		return errors.New("invalid assignee")
	}

	task, err := h.taskRepo.GetTaskByExternalID(message.TaskID)
	if err != nil {
		return err
	}

	h.accountRepo.Deposit(message.UserID, task.Reward);

	log.Printf("Task %d completed by %s", message.TaskID, message.UserID)

	return nil
}

func (h *Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountRepo.GetAllAccountsWithBalance()
	if err != nil {
		return
	}

	var accountsRes []accountsResponse
	for _, t := range accounts {
		accountsRes = append(accountsRes, accountsResponse{
			Username: t.Username,
			Balance:  t.Balance,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(getAccountsResponse{Accounts: accountsRes}); err != nil {
		log.Printf("Failed to encode accounts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}