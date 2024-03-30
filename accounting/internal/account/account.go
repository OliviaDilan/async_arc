package account

type Account struct {
	Username string
	Balance  int
}

type Repository interface {
	Create(username string) error
	GetAllAccountsWithBalance() ([]*Account, error)
	Deposit(username string, amount int) error
	Withdraw(username string, amount int) error
	GetAccountBalance(username string) (int, error)
}