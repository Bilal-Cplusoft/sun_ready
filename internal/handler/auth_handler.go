package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email     string `json:"email" example:"user@example.com"`
	Password  string `json:"password" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	CompanyID int    `json:"company_id" example:"1"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type AuthResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  interface{} `json:"user"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.CompanyID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.CompanyID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondJSON(w, http.StatusCreated, AuthResponse{Token: token, User: user})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, AuthResponse{Token: token, User: user})
}
