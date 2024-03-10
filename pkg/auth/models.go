package auth

type User struct {
	Username string
	Role Role
}

type decodeTokenRequest struct {
	Token string `json:"token"`
}

type decodeTokenResponse struct {
	Username string `json:"username"`
	Role     Role `json:"role"`
}

type getUsersResponse struct {
	Users []*User `json:"users"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleDeveloper  Role = "developer"
)