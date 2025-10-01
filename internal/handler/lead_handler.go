package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/go-chi/chi/v5"
)

type LeadHandler struct {
	leadService *service.LeadService
}

func NewLeadHandler(leadService *service.LeadService) *LeadHandler {
	return &LeadHandler{leadService: leadService}
}

// CreateLeadRequest represents the request body for creating a lead
type CreateLeadRequest struct {
	CompanyID           int      `json:"company_id" example:"1"`
	CreatorID           int      `json:"creator_id" example:"1"`
	Latitude            float64  `json:"latitude" example:"37.7749"`
	Longitude           float64  `json:"longitude" example:"-122.4194"`
	Address             string   `json:"address" example:"123 Solar St, San Francisco, CA 94102"`
	Source              int      `json:"source" example:"0"`
	PromoCode           *string  `json:"promo_code,omitempty" example:"SOLAR2025"`
	Is2D                bool     `json:"is_2d" example:"false"`
	KwhUsage            float64  `json:"kwh_usage" example:"12000"`
	SystemSize          float64  `json:"system_size,omitempty" example:"10.5"`
	PanelCount          int      `json:"panel_count,omitempty" example:"30"`
	PanelID             *int     `json:"panel_id,omitempty" example:"1"`
	InverterID          *int     `json:"inverter_id,omitempty" example:"1"`
	UtilityID           *int     `json:"utility_id,omitempty" example:"1"`
	RoofMaterial        *int     `json:"roof_material,omitempty" example:"1"`
}

// UpdateLeadRequest represents the request body for updating a lead
type UpdateLeadRequest struct {
	State               *int     `json:"state,omitempty" example:"1"`
	Latitude            *float64 `json:"latitude,omitempty" example:"37.7749"`
	Longitude           *float64 `json:"longitude,omitempty" example:"-122.4194"`
	Address             *string  `json:"address,omitempty" example:"123 Solar St, San Francisco, CA 94102"`
	PromoCode           *string  `json:"promo_code,omitempty" example:"SOLAR2025"`
	KwhUsage            *float64 `json:"kwh_usage,omitempty" example:"12000"`
	SystemSize          *float64 `json:"system_size,omitempty" example:"10.5"`
	PanelCount          *int     `json:"panel_count,omitempty" example:"30"`
	PanelID             *int     `json:"panel_id,omitempty" example:"1"`
	InverterID          *int     `json:"inverter_id,omitempty" example:"1"`
	UtilityID           *int     `json:"utility_id,omitempty" example:"1"`
	RoofMaterial        *int     `json:"roof_material,omitempty" example:"1"`
	AnnualProduction    *float64 `json:"annual_production,omitempty" example:"13000"`
	InstallationDate    *string  `json:"installation_date,omitempty" example:"2025-10-15"`
}

// LeadResponse represents the response for lead operations
type LeadResponse struct {
	Lead *models.Lead `json:"lead"`
}

// LeadsResponse represents the response for listing leads
type LeadsResponse struct {
	Leads []*models.Lead `json:"leads"`
	Total int            `json:"total"`
}

// Create godoc
// @Summary Create a new lead
// @Description Create a new lead with the provided details
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateLeadRequest true "Lead details"
// @Success 201 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads [post]
func (h *LeadHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	lead := &models.Lead{
		CompanyID:    req.CompanyID,
		CreatorID:    req.CreatorID,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		Source:       req.Source,
		PromoCode:    req.PromoCode,
		Is2D:         req.Is2D,
		KwhUsage:     req.KwhUsage,
		SystemSize:   req.SystemSize,
		PanelCount:   req.PanelCount,
		PanelID:      req.PanelID,
		InverterID:   req.InverterID,
		UtilityID:    req.UtilityID,
		RoofMaterial: req.RoofMaterial,
		State:        0, // Initialize as Progress
	}

	if err := h.leadService.Create(r.Context(), lead); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, LeadResponse{Lead: lead})
}

// GetByID godoc
// @Summary Get lead by ID
// @Description Get a lead by its ID
// @Tags leads
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lead ID"
// @Success 200 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /leads/{id} [get]
func (h *LeadHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Lead not found")
		return
	}

	respondJSON(w, http.StatusOK, LeadResponse{Lead: lead})
}

// Update godoc
// @Summary Update lead
// @Description Update a lead's details
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lead ID"
// @Param request body UpdateLeadRequest true "Lead update details"
// @Success 200 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads/{id} [put]
func (h *LeadHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Lead not found")
		return
	}

	var req UpdateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update only provided fields
	if req.State != nil {
		lead.State = *req.State
	}
	if req.Latitude != nil {
		lead.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		lead.Longitude = *req.Longitude
	}
	if req.Address != nil {
		lead.Address = *req.Address
	}
	if req.PromoCode != nil {
		lead.PromoCode = req.PromoCode
	}
	if req.KwhUsage != nil {
		lead.KwhUsage = *req.KwhUsage
	}
	if req.SystemSize != nil {
		lead.SystemSize = *req.SystemSize
	}
	if req.PanelCount != nil {
		lead.PanelCount = *req.PanelCount
	}
	if req.PanelID != nil {
		lead.PanelID = req.PanelID
	}
	if req.InverterID != nil {
		lead.InverterID = req.InverterID
	}
	if req.UtilityID != nil {
		lead.UtilityID = req.UtilityID
	}
	if req.RoofMaterial != nil {
		lead.RoofMaterial = req.RoofMaterial
	}
	if req.AnnualProduction != nil {
		lead.AnnualProduction = *req.AnnualProduction
	}
	if req.InstallationDate != nil {
		lead.InstallationDate = req.InstallationDate
	}

	if err := h.leadService.Update(r.Context(), lead); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, LeadResponse{Lead: lead})
}

// Delete godoc
// @Summary Delete lead
// @Description Delete a lead by ID
// @Tags leads
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lead ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads/{id} [delete]
func (h *LeadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	if err := h.leadService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete lead")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// List godoc
// @Summary List all leads
// @Description Get a list of all leads with pagination
// @Tags leads
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} LeadsResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads [get]
func (h *LeadHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	leads, err := h.leadService.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch leads")
		return
	}

	respondJSON(w, http.StatusOK, LeadsResponse{
		Leads: leads,
		Total: len(leads),
	})
}

// ListByCompany godoc
// @Summary List leads by company
// @Description Get a list of leads for a specific company
// @Tags leads
// @Produce json
// @Security BearerAuth
// @Param company_id path int true "Company ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} LeadsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads/company/{company_id} [get]
func (h *LeadHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyIDStr := chi.URLParam(r, "company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	leads, err := h.leadService.ListByCompany(r.Context(), companyID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch leads")
		return
	}

	respondJSON(w, http.StatusOK, LeadsResponse{
		Leads: leads,
		Total: len(leads),
	})
}

// UpdateState godoc
// @Summary Update lead state
// @Description Update the state of a lead
// @Tags leads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Lead ID"
// @Param state body map[string]int true "State"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /leads/{id}/state [put]
func (h *LeadHandler) UpdateState(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	var req map[string]int
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	state, ok := req["state"]
	if !ok {
		respondError(w, http.StatusBadRequest, "State is required")
		return
	}

	if err := h.leadService.UpdateState(r.Context(), id, state); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update lead state")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}
