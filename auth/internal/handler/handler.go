package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/OliviaDilan/async_arc/auth/internal/jwt"
	"github.com/OliviaDilan/async_arc/auth/internal/user"
)

type Handler struct {
	userRepo user.Repository
	jwt      jwt.Service
}

func NewHandler(userRepo user.Repository, jwt jwt.Service) *Handler {
	return &Handler{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {

	var req registerRequest

	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" || req.Role == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := user.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.userRepo.Create(&user); err != nil {
		log.Printf("Failed to create user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.userRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var usersRes []userResponce
	for _, u := range users {
		usersRes = append(usersRes, userResponce{
			Username: u.Username,
			Role:     u.Role,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(getUsersResponce{Users: usersRes}); err != nil {
		log.Printf("Failed to encode users: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req loginRequest
	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.Get(req.Username, req.Password)
	if err != nil {
		log.Printf("Failed to get user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.jwt.CreateToken(user)
	if err != nil {
		log.Printf("Failed to create token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loginResponce{Token: token}); err != nil {
		log.Printf("Failed to encode token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DecodeToken(w http.ResponseWriter, r *http.Request) {
	
	var req DecodeTokenRequest
	json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if req.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, err := h.jwt.DecodeToken(req.Token)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("Failed to decode token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(DecodeTokenResponce{
		Username: claims.Username,
		Role:     claims.Role,
	}); err != nil {
		log.Printf("Failed to encode token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}