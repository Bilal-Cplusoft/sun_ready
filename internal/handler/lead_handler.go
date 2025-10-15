package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/Bilal-Cplusoft/sun_ready/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type LeadHandler struct {
	leadRepo          *repo.LeadRepo
	lightFusionClient *client.LightFusionClient
	leadService *service.LeadService
	userRepo *repo.UserRepo
}

func NewLeadHandler(leadRepo *repo.LeadRepo, lightFusionClient *client.LightFusionClient, leadService *service.LeadService, userRepo *repo.UserRepo) *LeadHandler {
	return &LeadHandler{
		leadRepo:          leadRepo,
		lightFusionClient: lightFusionClient,
		leadService: leadService,
		userRepo: userRepo,
	}
}



// CreateLead godoc
// @Summary Create a new lead
// @Description Creates a new lead and optionally initiates 3D model generation
// @Tags leads
// @Accept json
// @Produce json
// @Param request body service.CreateLead true "Lead details"
// @Success 201 {object} LeadResponse "Returns lead and house IDs"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leads [post]
func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
    var req service.CreateLead
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, ok := middleware.GetUserID(r.Context())
    if !ok {
		http.Error(w, "User ID missing from context", http.StatusUnauthorized)
		return
    }
    companyID, _ := middleware.GetCompanyID(r.Context())
    user,_ := h.userRepo.GetByID(r.Context(),userID)

    effectiveCompanyID, err := h.userRepo.GetEffectiveCompanyID(r.Context(), user, companyID)
	if err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}
	response, err := h.leadService.CreateLead(r.Context(), req, userID, effectiveCompanyID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	leadResponse := LeadResponse{
		Success: response.Success,
		LeadID:  response.LeadID,
		HouseID: response.HouseID,
	}
	respondJSON(w, http.StatusCreated, leadResponse)
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
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	respondJSON(w, http.StatusOK, lead)
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
	var companyID, creatorID *int
	limit := 20
	offset := 0

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


	var leads []*models.Lead
	var total int64
	var err error

	leads, total, err = h.leadRepo.List(r.Context(), companyID, creatorID, limit, offset)
	if err != nil {
		log.Printf("Failed to list leads: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list leads")
		return
	}

	result := map[string]interface{}{
		"leads": leads,
		"total": total,
		"limit": limit,
		"offset": offset,
	}

	respondJSON(w, http.StatusOK, result)
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
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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
		respondError(w, http.StatusInternalServerError, "Failed to update lead")
		return
	}

	respondJSON(w, http.StatusOK, "success")
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
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	if err := h.leadRepo.Delete(r.Context(), id); err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to delete lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete lead")
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
		respondError(w, http.StatusBadRequest, "Invalid lead ID")
		return
	}

	lead, err := h.leadRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == models.ErrLeadNotFound {
			respondError(w, http.StatusNotFound, "Lead not found")
			return
		}
		log.Printf("Failed to get lead: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get lead")
		return
	}

	if !lead.Has3DModel() {
		respondError(w, http.StatusBadRequest, "Lead does not have a 3D model")
		return
	}

	// Get status from LightFusion
	status, err := h.lightFusionClient.GetProjectStatus(r.Context(), *lead.LightFusion3DProjectID, *lead.LightFusion3DHouseID)
	if err != nil {
		log.Printf("Failed to get 3D model status from LightFusion: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to sync 3D model status")
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
			respondError(w, http.StatusInternalServerError, "Failed to update lead")
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{
      "message": "synced to internal lead",
    })

}
