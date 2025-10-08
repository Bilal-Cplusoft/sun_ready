package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/go-chi/chi/v5"
)

type CompanyHandler struct {
	companyService *service.CompanyService
	userService    *service.UserService
}

func NewCompanyHandler(companyService *service.CompanyService, userService *service.UserService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		userService:    userService,
	}
}

// CreateCompanyRequest represents the request body for creating a company
type CreateCompanyRequest struct {
	Name                   string   `json:"name" example:"Acme Corp"`
	DisplayName            string   `json:"display_name" example:"Acme Corporation"`
	Description            string   `json:"description" example:"Leading solar company"`
	Code                   string   `json:"code" example:"ACME"`
	Slug                   string   `json:"slug" example:"acme-corp"`
	LogoPath               *string  `json:"logo_path" example:"https://example.com/logo.png"`
	SalesCommissionMin     *float64 `json:"sales_commission_min" example:"0.05"`
	SalesCommissionMax     *float64 `json:"sales_commission_max" example:"0.15"`
	SalesCommissionDefault *float64 `json:"sales_commission_default" example:"0.10"`
	BaselineAdder          *float64 `json:"baseline_adder" example:"100.00"`
}

// UpdateCompanyRequest represents the request body for updating a company
type UpdateCompanyRequest struct {
	Name                       *string  `json:"name,omitempty" example:"Acme Corp"`
	DisplayName                *string  `json:"display_name,omitempty" example:"Acme Corporation"`
	Description                *string  `json:"description,omitempty" example:"Leading solar company"`
	Code                       *string  `json:"code,omitempty" example:"ACME"`
	Slug                       *string  `json:"slug,omitempty" example:"acme-corp"`
	LogoPath                   *string  `json:"logo_path,omitempty" example:"https://example.com/logo.png"`
	SalesCommissionMin         *float64 `json:"sales_commission_min,omitempty" example:"0.05"`
	SalesCommissionMax         *float64 `json:"sales_commission_max,omitempty" example:"0.15"`
	SalesCommissionDefault     *float64 `json:"sales_commission_default,omitempty" example:"0.10"`
	Baseline                   *float64 `json:"baseline,omitempty" example:"1000.00"`
	BaselineAdder              *float64 `json:"baseline_adder,omitempty" example:"100.00"`
	BaselineAdderPctSalesComms *int     `json:"baseline_adder_pct_sales_comms,omitempty" example:"10"`
	ContractTag                *string  `json:"contract_tag,omitempty" example:"STANDARD"`
	IsActive                   *bool    `json:"is_active,omitempty" example:"true"`
}

// CompanyResponse represents the response for company operations
type CompanyResponse struct {
	Company *models.Company `json:"company"`
}

// CompaniesResponse represents the response for listing companies
type CompaniesResponse struct {
	Companies []*models.Company `json:"companies"`
	Total     int               `json:"total"`
}

// Create godoc
// @Summary Create a new company
// @Description Create a new company with the provided details
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCompanyRequest true "Company details"
// @Success 201 {object} CompanyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	company := &models.Company{
		Name:                   req.Name,
		DisplayName:            req.DisplayName,
		Description:            req.Description,
		Code:                   req.Code,
		Slug:                   req.Slug,
		LogoPath:               req.LogoPath,
		SalesCommissionMin:     req.SalesCommissionMin,
		SalesCommissionMax:     req.SalesCommissionMax,
		SalesCommissionDefault: req.SalesCommissionDefault,
		BaselineAdder:          req.BaselineAdder,
		IsActive:               true,
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

	respondJSON(w, http.StatusCreated, CompanyResponse{Company: company})
}

// GetByID godoc
// @Summary Get company by ID
// @Description Get a company by its ID
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} CompanyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/companies/{id} [get]
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.companyService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Company not found")
		return
	}

	respondJSON(w, http.StatusOK, CompanyResponse{Company: company})
}

// GetBySlug godoc
// @Summary Get company by slug
// @Description Get a company by its slug
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Company Slug"
// @Success 200 {object} CompanyResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/companies/slug/{slug} [get]
func (h *CompanyHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	company, err := h.companyService.GetBySlug(r.Context(), slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "Company not found")
		return
	}

	respondJSON(w, http.StatusOK, CompanyResponse{Company: company})
}

// Update godoc
// @Summary Update company
// @Description Update a company's details
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body UpdateCompanyRequest true "Company update details"
// @Success 200 {object} CompanyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies/{id} [put]
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	company, err := h.companyService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Company not found")
		return
	}

	var req UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update only provided fields
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.DisplayName != nil {
		company.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		company.Description = *req.Description
	}
	if req.Code != nil {
		company.Code = *req.Code
	}
	if req.Slug != nil {
		company.Slug = *req.Slug
	}
	if req.LogoPath != nil {
		company.LogoPath = req.LogoPath
	}
	if req.SalesCommissionMin != nil {
		company.SalesCommissionMin = req.SalesCommissionMin
	}
	if req.SalesCommissionMax != nil {
		company.SalesCommissionMax = req.SalesCommissionMax
	}
	if req.SalesCommissionDefault != nil {
		company.SalesCommissionDefault = req.SalesCommissionDefault
	}
	if req.Baseline != nil {
		company.Baseline = req.Baseline
	}
	if req.BaselineAdder != nil {
		company.BaselineAdder = req.BaselineAdder
	}
	if req.BaselineAdderPctSalesComms != nil {
		company.BaselineAdderPctSalesComms = req.BaselineAdderPctSalesComms
	}
	if req.ContractTag != nil {
		company.ContractTag = req.ContractTag
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}

	company.Sanitize()
	if err := company.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.companyService.Update(r.Context(), company); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update company")
		return
	}

	respondJSON(w, http.StatusOK, CompanyResponse{Company: company})
}

// Delete godoc
// @Summary Delete company
// @Description Delete a company by ID
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies/{id} [delete]
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	if err := h.companyService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete company")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// List godoc
// @Summary List all companies
// @Description Get a list of all companies with optional filtering
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param referred_by_user_id query int false "Filter by referred by user ID"
// @Success 200 {object} CompaniesResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies [get]
func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	referredByStr := r.URL.Query().Get("referred_by_user_id")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	companies, err := h.companyService.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	// Filter by referred_by_user_id if provided
	if referredByStr != "" {
		if referredByID, err := strconv.Atoi(referredByStr); err == nil {
			var filtered []*models.Company
			for _, c := range companies {
				if c.ReferredByUserID != nil && *c.ReferredByUserID == referredByID {
					filtered = append(filtered, c)
				}
			}
			companies = filtered
		}
	}

	respondJSON(w, http.StatusOK, CompaniesResponse{
		Companies: companies,
		Total:     len(companies),
	})
}

// FindAll godoc
// @Summary Get all companies
// @Description Get all companies without pagination with optional filtering
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param referred_by_user_id query int false "Filter by referred by user ID"
// @Success 200 {object} CompaniesResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/companies/all [get]
func (h *CompanyHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	referredByStr := r.URL.Query().Get("referred_by_user_id")

	companies, err := h.companyService.FindAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch companies")
		return
	}

	// Filter by referred_by_user_id if provided
	if referredByStr != "" {
		if referredByID, err := strconv.Atoi(referredByStr); err == nil {
			var filtered []*models.Company
			for _, c := range companies {
				if c.ReferredByUserID != nil && *c.ReferredByUserID == referredByID {
					filtered = append(filtered, c)
				}
			}
			companies = filtered
		}
	}

	respondJSON(w, http.StatusOK, CompaniesResponse{
		Companies: companies,
		Total:     len(companies),
	})
}
