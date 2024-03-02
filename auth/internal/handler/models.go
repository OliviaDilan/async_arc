package handler

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponce struct {
	Token string `json:"token"`
}

type getUsersResponce struct {
	Users []userResponce `json:"users"`
}

type userResponce struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type DecodeTokenRequest struct {
	Token string `json:"token"`
}

type DecodeTokenResponce struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}