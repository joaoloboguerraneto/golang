package handler

import (
	"encoding/json"
	"net/http"

	"github.com/workshop/golang-login/internal/models"
	"github.com/workshop/golang-login/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "JSON inválido"})
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "Campos obrigatórios"})
		return
	}

	err := h.authService.RegisterUser(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, APIResponse{Success: true, Message: "✅ Usuário cadastrado!"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Message: "JSON inválido"})
		return
	}

	result, err := h.authService.Login(r.Context(), models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		status := http.StatusUnauthorized
		if err == service.ErrUserNotFound {
			status = http.StatusNotFound
		}
		respondJSON(w, status, APIResponse{Success: false, Message: result.Message})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: result.Message, Data: result.User})
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.ListUsers(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Message: err.Error()})
		return
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: users})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
