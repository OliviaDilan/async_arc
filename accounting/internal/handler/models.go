package handler

type accountsResponse struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

type getAccountsResponse struct {
	Accounts []accountsResponse `json:"accounts"`
}