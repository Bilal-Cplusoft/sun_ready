package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user.ID = id
	if err := h.userService.Update(r.Context(), &user); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, err := h.userService.List(r.Context(), companyID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
