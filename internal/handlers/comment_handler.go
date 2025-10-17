package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"goida/internal/middleware"
	"goida/internal/models"
	"goida/internal/services"
)

type CommentHandler struct {
	service   services.CommentService
	validator *middleware.Validator
}

func NewCommentHandler(service services.CommentService, validator *middleware.Validator) *CommentHandler {
	return &CommentHandler{service: service, validator: validator}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validator.ValidateStruct(&req); err != nil {
		validationErrors := h.validator.FormatValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Validation failed", "details": validationErrors})
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User context not found", http.StatusInternalServerError)
		return
	}

	comment, err := h.service.Create(r.Context(), articleID, claims.UserID, &req)
	if err != nil {
		logrus.Errorf("Failed to create comment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	limit, offset := 10, 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	items, err := h.service.ListByArticle(r.Context(), articleID, limit, offset)
	if err != nil {
		logrus.Errorf("Failed to list comments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validator.ValidateStruct(&req); err != nil {
		validationErrors := h.validator.FormatValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Validation failed", "details": validationErrors})
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User context not found", http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdateOwned(r.Context(), id64, claims.UserID, &req); err != nil {
		if err.Error() == "not found or not owner" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User context not found", http.StatusInternalServerError)
		return
	}

	if err := h.service.DeleteOwned(r.Context(), id64, claims.UserID); err != nil {
		if err.Error() == "not found or not owner" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
