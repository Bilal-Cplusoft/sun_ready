package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DealHandler struct {
	dealService *service.DealService
}

func NewDealHandler(dealService *service.DealService) *DealHandler {
	return &DealHandler{dealService: dealService}
}


type CreateDealRequest struct {
	ProjectID           int     `json:"project_id" example:"1"`
	SystemSize          float64 `json:"system_size" example:"10.5"`
	PanelCount          int     `json:"panel_count" example:"30"`
	SalesID             int     `json:"sales_id" example:"1"`
	HomeownerID         int     `json:"homeowner_id" example:"2"`
	DocumentID          *int    `json:"document_id,omitempty" example:"1"`
	FinancingOptionID   *int    `json:"financing_option_id,omitempty" example:"1"`
	FinancingProvider   string  `json:"financing_provider,omitempty" example:"SunPower Financial"`
	TargetEPC           float64 `json:"target_epc" example:"2.50"`
	TotalCost           float64 `json:"total_cost" example:"25000.00"`
	HardwareCost        float64 `json:"hardware_cost" example:"15000.00"`
	InstallationCost    float64 `json:"installation_cost" example:"8000.00"`
	SalesCommissionCost float64 `json:"sales_commission_cost" example:"2000.00"`
	Profit              float64 `json:"profit" example:"5000.00"`
	CompanyID           int     `json:"company_id" example:"1"`
	Address             string  `json:"address,omitempty" example:"123 Solar Street, CA 90210"`
	ConsumptionKWH      int     `json:"consumption_kwh,omitempty" example:"12000"`
	ProductionKWH       int     `json:"production_kwh,omitempty" example:"13000"`
}


type UpdateDealRequest struct {
	SystemSize          *float64    `json:"system_size,omitempty" example:"10.5"`
	PanelCount          *int        `json:"panel_count,omitempty" example:"30"`
	DocumentID          *int        `json:"document_id,omitempty" example:"1"`
	FinancingOptionID   *int        `json:"financing_option_id,omitempty" example:"1"`
	FinancingProvider   *string     `json:"financing_provider,omitempty" example:"SunPower Financial"`
	TargetEPC           *float64    `json:"target_epc,omitempty" example:"2.50"`
	TotalCost           *float64    `json:"total_cost,omitempty" example:"25000.00"`
	HardwareCost        *float64    `json:"hardware_cost,omitempty" example:"15000.00"`
	InstallationCost    *float64    `json:"installation_cost,omitempty" example:"8000.00"`
	SalesCommissionCost *float64    `json:"sales_commission_cost,omitempty" example:"2000.00"`
	Profit              *float64    `json:"profit,omitempty" example:"5000.00"`
	Status              *string     `json:"status,omitempty" example:"approved"`
	SignedAt            *time.Time  `json:"signed_at,omitempty" example:"2025-10-01T10:00:00Z"`
	ApprovedAt          *time.Time  `json:"approved_at,omitempty" example:"2025-10-02T10:00:00Z"`
	InstalledAt         *time.Time  `json:"installed_at,omitempty" example:"2025-10-15T10:00:00Z"`
	Address             *string     `json:"address,omitempty" example:"123 Solar Street, CA 90210"`
	ConsumptionKWH      *int        `json:"consumption_kwh,omitempty" example:"12000"`
	ProductionKWH       *int        `json:"production_kwh,omitempty" example:"13000"`
}

// DealResponse represents the response for deal operations
type DealResponse struct {
	Deal *models.Deal `json:"deal"`
}

// DealsResponse represents the response for listing deals
type DealsResponse struct {
	Deals []*models.Deal `json:"deals"`
	Total int            `json:"total"`
}

// Create godoc
// @Summary Create a new deal
// @Description Create a new deal with the provided details
// @Tags deals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDealRequest true "Deal details"
// @Success 201 {object} DealResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals [post]
func (h *DealHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateDealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	deal := &models.Deal{
		UUID:                uuid.New().String(),
		ProjectID:           req.ProjectID,
		SystemSize:          req.SystemSize,
		PanelCount:          req.PanelCount,
		SalesID:             req.SalesID,
		HomeownerID:         req.HomeownerID,
		DocumentID:          req.DocumentID,
		FinancingOptionID:   req.FinancingOptionID,
		FinancingProvider:   req.FinancingProvider,
		TargetEPC:           req.TargetEPC,
		TotalCost:           req.TotalCost,
		HardwareCost:        req.HardwareCost,
		InstallationCost:    req.InstallationCost,
		SalesCommissionCost: req.SalesCommissionCost,
		Profit:              req.Profit,
		CompanyID:           req.CompanyID,
		Address:             req.Address,
		ConsumptionKWH:      req.ConsumptionKWH,
		ProductionKWH:       req.ProductionKWH,
		Status:              "pending",
		Archive:             false,
	}

	if err := h.dealService.Create(r.Context(), deal); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, DealResponse{Deal: deal})
}

// GetByID godoc
// @Summary Get deal by ID
// @Description Get a deal by its ID
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Success 200 {object} DealResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/deals/{id} [get]
func (h *DealHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	deal, err := h.dealService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Deal not found")
		return
	}

	respondJSON(w, http.StatusOK, DealResponse{Deal: deal})
}

// GetByUUID godoc
// @Summary Get deal by UUID
// @Description Get a deal by its UUID
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param uuid path string true "Deal UUID"
// @Success 200 {object} DealResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/deals/uuid/{uuid} [get]
func (h *DealHandler) GetByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	deal, err := h.dealService.GetByUUID(r.Context(), uuid)
	if err != nil {
		respondError(w, http.StatusNotFound, "Deal not found")
		return
	}

	respondJSON(w, http.StatusOK, DealResponse{Deal: deal})
}

// Update godoc
// @Summary Update deal
// @Description Update a deal's details
// @Tags deals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Param request body UpdateDealRequest true "Deal update details"
// @Success 200 {object} DealResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/{id} [put]
func (h *DealHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	deal, err := h.dealService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Deal not found")
		return
	}

	var req UpdateDealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update only provided fields
	if req.SystemSize != nil {
		deal.SystemSize = *req.SystemSize
	}
	if req.PanelCount != nil {
		deal.PanelCount = *req.PanelCount
	}
	if req.DocumentID != nil {
		deal.DocumentID = req.DocumentID
	}
	if req.FinancingOptionID != nil {
		deal.FinancingOptionID = req.FinancingOptionID
	}
	if req.FinancingProvider != nil {
		deal.FinancingProvider = *req.FinancingProvider
	}
	if req.TargetEPC != nil {
		deal.TargetEPC = *req.TargetEPC
	}
	if req.TotalCost != nil {
		deal.TotalCost = *req.TotalCost
	}
	if req.HardwareCost != nil {
		deal.HardwareCost = *req.HardwareCost
	}
	if req.InstallationCost != nil {
		deal.InstallationCost = *req.InstallationCost
	}
	if req.SalesCommissionCost != nil {
		deal.SalesCommissionCost = *req.SalesCommissionCost
	}
	if req.Profit != nil {
		deal.Profit = *req.Profit
	}
	if req.Status != nil {
		deal.Status = *req.Status
	}
	if req.SignedAt != nil {
		deal.SignedAt = req.SignedAt
	}
	if req.ApprovedAt != nil {
		deal.ApprovedAt = req.ApprovedAt
	}
	if req.InstalledAt != nil {
		deal.InstalledAt = req.InstalledAt
	}
	if req.Address != nil {
		deal.Address = *req.Address
	}
	if req.ConsumptionKWH != nil {
		deal.ConsumptionKWH = *req.ConsumptionKWH
	}
	if req.ProductionKWH != nil {
		deal.ProductionKWH = *req.ProductionKWH
	}

	if err := h.dealService.Update(r.Context(), deal); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, DealResponse{Deal: deal})
}

// Delete godoc
// @Summary Delete deal
// @Description Delete a deal by ID
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/{id} [delete]
func (h *DealHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	if err := h.dealService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete deal")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// List godoc
// @Summary List all deals
// @Description Get a list of all deals with pagination
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} DealsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/ [get]
func (h *DealHandler) List(w http.ResponseWriter, r *http.Request) {
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

	deals, err := h.dealService.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch deals")
		return
	}

	respondJSON(w, http.StatusOK, DealsResponse{
		Deals: deals,
		Total: len(deals),
	})
}

// ListByCompany godoc
// @Summary List deals by company
// @Description Get a list of deals for a specific company
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param company_id path int true "Company ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} DealsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/company/{company_id} [get]
func (h *DealHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
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

	deals, err := h.dealService.ListByCompany(r.Context(), companyID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch deals")
		return
	}

	respondJSON(w, http.StatusOK, DealsResponse{
		Deals: deals,
		Total: len(deals),
	})
}

// ListSigned godoc
// @Summary List signed deals
// @Description Get a list of signed deals for a company
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param company_id path int true "Company ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} DealsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/company/{company_id}/signed [get]
func (h *DealHandler) ListSigned(w http.ResponseWriter, r *http.Request) {
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

	deals, err := h.dealService.ListSigned(r.Context(), companyID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch signed deals")
		return
	}

	respondJSON(w, http.StatusOK, DealsResponse{
		Deals: deals,
		Total: len(deals),
	})
}

// Archive godoc
// @Summary Archive a deal
// @Description Archive a deal by ID
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/{id}/archive [post]
func (h *DealHandler) Archive(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	if err := h.dealService.Archive(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to archive deal")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Unarchive godoc
// @Summary Unarchive a deal
// @Description Unarchive a deal by ID
// @Tags deals
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/deals/{id}/unarchive [post]
func (h *DealHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	if err := h.dealService.Unarchive(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to unarchive deal")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}
