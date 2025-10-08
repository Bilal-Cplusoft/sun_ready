package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
)

// AddCompanyRequest represents the request for adding a company with user migration
type AddCompanyRequest struct {
	CompanyName string `json:"company_name" example:"New Solar Company"`
	MainUserID  int    `json:"main_user_id" example:"1"`
}

// AddCompany godoc
// @Summary Add company with user migration
// @Description Creates a new company and migrates the main user and their descendants to it
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddCompanyRequest true "Company and user details"
// @Success 201 {object} CompanyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies/add [post]
func (h *CompanyHandler) AddCompany(w http.ResponseWriter, r *http.Request) {
	var req AddCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.CompanyName == "" {
		respondError(w, http.StatusBadRequest, "Company name is required")
		return
	}

	if req.MainUserID == 0 {
		respondError(w, http.StatusBadRequest, "Main user ID is required")
		return
	}

	// Get the main user
	mainUser, err := h.userService.GetByID(r.Context(), req.MainUserID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Main user not found")
		return
	}

	// Store the creator ID before clearing it
	var referredByUserID *int
	if mainUser.CreatorID != nil {
		referredByUserID = mainUser.CreatorID
	}

	// Create the company
	company := &models.Company{
		Name:             req.CompanyName,
		Slug:             generateSlug(req.CompanyName),
		DisplayName:      req.CompanyName,
		IsActive:         true,
		ReferredByUserID: referredByUserID,
	}

	company.Sanitize()
	if err := company.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.companyService.Create(r.Context(), company); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create company")
		return
	}

	// Update main user's company
	mainUser.CreatorID = nil
	mainUser.CompanyID = company.ID
	if err := h.userService.Update(r.Context(), mainUser); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update main user")
		return
	}

	// Get all descendants of the main user
	descendantIDs, err := h.userService.GetDescendantIDs(r.Context(), req.MainUserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get descendant users")
		return
	}

	// Update all descendants to the new company
	if len(descendantIDs) > 0 {
		descendants, err := h.userService.FindByIDs(r.Context(), descendantIDs)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to fetch descendant users")
			return
		}

		for _, user := range descendants {
			user.CompanyID = company.ID
			if err := h.userService.Update(r.Context(), user); err != nil {
				respondError(w, http.StatusInternalServerError, "Failed to update descendant user")
				return
			}
		}
	}

	// Update company's referred_by_user_id
	if referredByUserID != nil {
		company.ReferredByUserID = referredByUserID
		if err := h.companyService.Update(r.Context(), company); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to update company referral")
			return
		}
	}

	respondJSON(w, http.StatusCreated, CompanyResponse{Company: company})
}

// generateSlug generates a URL-friendly slug from a company name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters
	var result []rune
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	return string(result)
}
