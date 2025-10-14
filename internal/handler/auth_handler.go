package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
)

type AuthHandler struct {
	authService *service.AuthService
	sendGridClient *client.SendGridClient
}

func NewAuthHandler(authService *service.AuthService,sendGridClient *client.SendGridClient) *AuthHandler {
	return &AuthHandler{authService: authService, sendGridClient: sendGridClient}
}

type RegisterRequest struct {
	Email     string `json:"email" example:"user@example.com"`
	Password  string `json:"password" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	CompanyID int    `json:"company_id" example:"1"`
	Address   string `json:"address" example:"123 Main St, Anytown, USA"`
	Phone     string `json:"phone" example:"555-123-4567"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type AuthResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  any `json:"user"`
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
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.CompanyID,req.Address,req.Phone)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.CompanyID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	err = h.sendGridClient.SendWelcomeEmail(req.Email,req.FirstName)
	if err != nil {
		respondError(w, http.StatusConflict, "Error sending welcome email")
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
// @Router /api/auth/login [post]
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
