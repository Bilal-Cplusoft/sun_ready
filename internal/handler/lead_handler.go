package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"github.com/go-chi/chi/v5"
)

type LeadHandler struct {
	leadRepo          *repo.LeadRepo
	lightFusionClient *client.LightFusionClient
}

func NewLeadHandler(leadRepo *repo.LeadRepo, lightFusionClient *client.LightFusionClient) *LeadHandler {
	return &LeadHandler{
		leadRepo:          leadRepo,
		lightFusionClient: lightFusionClient,
	}
}

// CreateLeadRequest represents the request to create a new lead
type CreateLeadRequest struct {
	CompanyID           int     `json:"company_id" example:"1"`
	CreatorID           int     `json:"creator_id" example:"1"`
	Latitude            float64 `json:"latitude" example:"37.7749"`
	Longitude           float64 `json:"longitude" example:"-122.4194"`
	Address             string  `json:"address" example:"123 Solar St, San Francisco, CA 94102"`
	Source              int     `json:"source" example:"2"`
	PromoCode           *string `json:"promo_code,omitempty" example:"SOLAR2025"`
	KwhUsage            float64 `json:"kwh_usage" example:"12000"`
	SystemSize          float64 `json:"system_size" example:"10.5"`
	PanelCount          int     `json:"panel_count" example:"30"`
	Create3DModel       bool    `json:"create_3d_model" example:"false"`
}

// LeadResponse represents a lead in API responses
type LeadResponse struct {
	ID                     int     `json:"id" example:"1"`
	CreatedAt              string  `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt              string  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
	CompanyID              int     `json:"company_id" example:"1"`
	CreatorID              int     `json:"creator_id" example:"1"`
	Latitude               float64 `json:"latitude" example:"37.7749"`
	Longitude              float64 `json:"longitude" example:"-122.4194"`
	Address                string  `json:"address" example:"123 Solar St, San Francisco, CA 94102"`
	State                  int     `json:"state" example:"0"`
	Source                 int     `json:"source" example:"2"`
	PromoCode              *string `json:"promo_code,omitempty" example:"SOLAR2025"`
	KwhUsage               float64 `json:"kwh_usage" example:"12000"`
	SystemSize             float64 `json:"system_size" example:"10.5"`
	PanelCount             int     `json:"panel_count" example:"30"`
	AnnualProduction       float64 `json:"annual_production" example:"13000"`
	ExternalLeadID         *int    `json:"external_lead_id,omitempty" example:"12345"`
	SyncStatus             string  `json:"sync_status" example:"synced"`
	Has3DModel             bool    `json:"has_3d_model" example:"true"`
	LightFusion3DProjectID *int    `json:"lightfusion_3d_project_id,omitempty" example:"123"`
	LightFusion3DHouseID   *int    `json:"lightfusion_3d_house_id,omitempty" example:"456"`
	Model3DStatus          *string `json:"model_3d_status,omitempty" example:"completed"`
}

// CreateLead godoc
// @Summary Create a new lead
// @Description Creates a new lead and optionally initiates 3D model generation
// @Tags leads
// @Accept json
// @Produce json
// @Param request body CreateLeadRequest true "Lead details"
// @Success 201 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads [post]
func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	var req CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create the lead
	lead := &models.Lead{
		CompanyID:        req.CompanyID,
		CreatorID:        req.CreatorID,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		Address:          req.Address,
		State:            int(models.LeadStateInitialized),
		Source:           req.Source,
		PromoCode:        req.PromoCode,
		KwhUsage:         req.KwhUsage,
		SystemSize:       req.SystemSize,
		PanelCount:       req.PanelCount,
	}

	if err := h.leadRepo.Create(r.Context(), lead); err != nil {
		log.Printf("Failed to create lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create lead")
		return
	}

	// Optionally create 3D model
	if req.Create3DModel && h.lightFusionClient != nil {
		// This would typically be done asynchronously
		log.Printf("3D model creation requested for lead %d - implement async processing", lead.ID)
	}

	response := h.leadToResponse(lead)
	respondWithJSON(w, http.StatusCreated, response)
}

// GetLead godoc
// @Summary Get a lead by ID
// @Description Retrieves a lead by its ID
// @Tags leads
// @Produce json
// @Param id path int true "Lead ID"
// @Success 200 {object} LeadResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/leads/{id} [get]
func (h *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondWithError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	response := h.leadToResponse(lead)
	respondWithJSON(w, http.StatusOK, response)
}

// ListLeads godoc
// @Summary List leads
// @Description Retrieves a paginated list of leads
// @Tags leads
// @Produce json
// @Param company_id query int false "Filter by company ID"
// @Param creator_id query int false "Filter by creator ID"
// @Param has_3d_model query bool false "Filter leads with 3D models"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /api/leads [get]
func (h *LeadHandler) ListLeads(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	var companyID, creatorID *int
	limit := 20
	offset := 0
	has3DModel := false

	if companyIDStr := r.URL.Query().Get("company_id"); companyIDStr != "" {
		if id, err := strconv.Atoi(companyIDStr); err == nil {
			companyID = &id
		}
	}

	if creatorIDStr := r.URL.Query().Get("creator_id"); creatorIDStr != "" {
		if id, err := strconv.Atoi(creatorIDStr); err == nil {
			creatorID = &id
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	if has3DModelStr := r.URL.Query().Get("has_3d_model"); has3DModelStr == "true" {
		has3DModel = true
	}

	var leads []*models.Lead
	var total int64
	var err error

	if has3DModel {
		leads, total, err = h.leadRepo.ListWith3DModels(r.Context(), companyID, limit, offset)
	} else {
		leads, total, err = h.leadRepo.List(r.Context(), companyID, creatorID, limit, offset)
	}

	if err != nil {
		log.Printf("Failed to list leads: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list leads")
		return
	}

	// Convert to response format
	responses := make([]LeadResponse, len(leads))
	for i, lead := range leads {
		responses[i] = h.leadToResponse(lead)
	}

	result := map[string]interface{}{
		"leads": responses,
		"total": total,
		"limit": limit,
		"offset": offset,
	}

	respondWithJSON(w, http.StatusOK, result)
}

// UpdateLead godoc
// @Summary Update a lead
// @Description Updates an existing lead
// @Tags leads
// @Accept json
// @Produce json
// @Param id path int true "Lead ID"
// @Param request body map[string]interface{} true "Lead updates"
// @Success 200 {object} LeadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads/{id} [put]
func (h *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondWithError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Apply updates (simplified - in production you'd want more validation)
	if address, ok := updates["address"].(string); ok {
		lead.Address = address
	}
	if kwhUsage, ok := updates["kwh_usage"].(float64); ok {
		lead.KwhUsage = kwhUsage
	}
	if systemSize, ok := updates["system_size"].(float64); ok {
		lead.SystemSize = systemSize
	}
	if panelCount, ok := updates["panel_count"].(float64); ok {
		lead.PanelCount = int(panelCount)
	}
	if annualProduction, ok := updates["annual_production"].(float64); ok {
		lead.AnnualProduction = annualProduction
	}

	if err := h.leadRepo.Update(r.Context(), lead); err != nil {
		log.Printf("Failed to update lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update lead")
		return
	}

	response := h.leadToResponse(lead)
	respondWithJSON(w, http.StatusOK, response)
}

// DeleteLead godoc
// @Summary Delete a lead
// @Description Deletes a lead by ID
// @Tags leads
// @Param id path int true "Lead ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads/{id} [delete]
func (h *LeadHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	if err := h.leadRepo.Delete(r.Context(), id); err != nil {
		if err == models.ErrLeadNotFound {
			respondWithError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to delete lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete lead")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SyncLead3DStatus godoc
// @Summary Sync 3D model status for a lead
// @Description Updates the 3D model status from LightFusion API
// @Tags leads
// @Param id path int true "Lead ID"
// @Success 200 {object} LeadResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads/{id}/sync-3d-status [post]
func (h *LeadHandler) SyncLead3DStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondWithError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	if !lead.Has3DModel() {
		respondWithError(w, http.StatusBadRequest, "Lead does not have a 3D model")
		return
	}

	// Get status from LightFusion
	status, err := h.lightFusionClient.GetProjectStatus(r.Context(), *lead.LightFusion3DProjectID, *lead.LightFusion3DHouseID)
	if err != nil {
		log.Printf("Failed to get 3D model status from LightFusion: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to sync 3D model status")
		return
	}

	// Update lead based on status
	if status.LeadCompletion != nil {
		leadData := status.LeadCompletion.Lead
		
		if leadData.State == 1 { // Done
			lead.Update3DModelStatus("completed")
		} else if leadData.State == 2 { // Errored
			lead.Update3DModelStatus("failed")
		} else {
			lead.Update3DModelStatus("processing")
		}

		// Update other fields if available
		if leadData.House.SystemSize > 0 {
			lead.SystemSize = float64(leadData.House.SystemSize)
		}
		if leadData.House.PanelCount > 0 {
			lead.PanelCount = leadData.House.PanelCount
		}
		if leadData.Production.Annual > 0 {
			lead.AnnualProduction = leadData.Production.Annual
		}

		if err := h.leadRepo.Update(r.Context(), lead); err != nil {
			log.Printf("Failed to update lead after sync: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to update lead")
			return
		}
	}

	response := h.leadToResponse(lead)
	respondWithJSON(w, http.StatusOK, response)
}

// Helper function to convert Lead model to response
func (h *LeadHandler) leadToResponse(lead *models.Lead) LeadResponse {
	return LeadResponse{
		ID:                     lead.ID,
		CreatedAt:              lead.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:              lead.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CompanyID:              lead.CompanyID,
		CreatorID:              lead.CreatorID,
		Latitude:               lead.Latitude,
		Longitude:              lead.Longitude,
		Address:                lead.Address,
		State:                  lead.State,
		Source:                 lead.Source,
		PromoCode:              lead.PromoCode,
		KwhUsage:               lead.KwhUsage,
		SystemSize:             lead.SystemSize,
		PanelCount:             lead.PanelCount,
		AnnualProduction:       lead.AnnualProduction,
		ExternalLeadID:         lead.ExternalLeadID,
		SyncStatus:             lead.SyncStatus,
		Has3DModel:             lead.Has3DModel(),
		LightFusion3DProjectID: lead.LightFusion3DProjectID,
		LightFusion3DHouseID:   lead.LightFusion3DHouseID,
		Model3DStatus:          lead.Model3DStatus,
	}
}

// SetupRoutes sets up the lead routes
func (h *LeadHandler) SetupRoutes(r chi.Router) {
	r.Route("/api/leads", func(r chi.Router) {
		r.Post("/", h.CreateLead)
		r.Get("/", h.ListLeads)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetLead)
			r.Put("/", h.UpdateLead)
			r.Delete("/", h.DeleteLead)
			r.Post("/sync-3d-status", h.SyncLead3DStatus)
		})
	})
}