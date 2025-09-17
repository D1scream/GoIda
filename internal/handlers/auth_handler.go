package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"goida/internal/middleware"
	"goida/internal/models"
	"goida/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *middleware.Validator
}

func NewAuthHandler(authService *services.AuthService, validator *middleware.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		validationErrors := h.validator.FormatValidationErrors(err)
		response := map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrors,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.authService.Authenticate(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User context not found", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role: &models.Role{
			Name: claims.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
