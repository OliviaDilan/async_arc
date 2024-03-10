package account

import (
	"fmt"
)

type inMemoryRepository struct {
	accounts map[string]*Account
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		accounts: make(map[string]*Account),
	}
}

func (r *inMemoryRepository) Create(username string) error {
	
	if _, ok := r.accounts[username]; ok {
		return fmt.Errorf("user already exists")
	}

	r.accounts[username] = &Account{
		Username: username,
		Balance:  0,
	}
	return nil
}

func (r *inMemoryRepository) GetAllAccountsWithBalance() ([]*Account, error) {
	var accounts []*Account
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *inMemoryRepository) Deposit(username string, amount int) error {
	account, ok := r.accounts[username]
	if !ok {
		return fmt.Errorf("user not found")
	}
	account.Balance += amount
	return nil
}

func (r *inMemoryRepository) Withdraw(username string, amount int) error {
	account, ok := r.accounts[username]
	if !ok {
		return fmt.Errorf("user not found")
	}
	account.Balance -= amount
	return nil
}

func (r *inMemoryRepository) GetAccountBalance(username string) (int, error) {
	account, ok := r.accounts[username]
	if !ok {
		return 0, fmt.Errorf("user not found")
	}
	return account.Balance, nil
}