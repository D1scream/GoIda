package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"goida/internal/middleware"
	"goida/internal/models"
	"goida/internal/repository"
)

type AuthCredentialsHandler struct {
	authCredentialsRepo repository.AuthCredentialsRepository
	validator           *middleware.Validator
}

func NewAuthCredentialsHandler(authCredentialsRepo repository.AuthCredentialsRepository, validator *middleware.Validator) *AuthCredentialsHandler {
	return &AuthCredentialsHandler{
		authCredentialsRepo: authCredentialsRepo,
		validator:           validator,
	}
}

func (h *AuthCredentialsHandler) CreateCredentials(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAuthCredentialsRequest
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

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	credentials := &models.AuthCredentials{
		UserID:   req.UserID,
		Login:    req.Login,
		Password: string(hashedPassword),
	}

	err = h.authCredentialsRepo.Create(credentials)
	if err != nil {
		logrus.Errorf("Failed to create auth credentials: %v", err)
		http.Error(w, "Failed to create credentials", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Credentials created successfully",
		"user_id": credentials.UserID,
		"login":   credentials.Login,
	})
}

func (h *AuthCredentialsHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	credentials, err := h.authCredentialsRepo.GetByUserID(userID)
	if err != nil {
		logrus.Errorf("Failed to get auth credentials: %v", err)
		http.Error(w, "Credentials not found", http.StatusNotFound)
		return
	}

	// Не возвращаем пароль
	response := map[string]interface{}{
		"user_id":    credentials.UserID,
		"login":      credentials.Login,
		"created_at": credentials.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthCredentialsHandler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateAuthCredentialsRequest
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

	// Получаем существующие учетные данные
	credentials, err := h.authCredentialsRepo.GetByUserID(userID)
	if err != nil {
		logrus.Errorf("Failed to get auth credentials: %v", err)
		http.Error(w, "Credentials not found", http.StatusNotFound)
		return
	}

	// Обновляем только переданные поля
	if req.Login != "" {
		credentials.Login = req.Login
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Errorf("Failed to hash password: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		credentials.Password = string(hashedPassword)
	}

	err = h.authCredentialsRepo.Update(userID, credentials)
	if err != nil {
		logrus.Errorf("Failed to update auth credentials: %v", err)
		http.Error(w, "Failed to update credentials", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Credentials updated successfully",
		"user_id": credentials.UserID,
		"login":   credentials.Login,
	})
}
