package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
)

type Project3DHandler struct {
	lightFusionClient *client.LightFusionClient
	leadRepo          *repo.LeadRepo
}

func NewProject3DHandler(lightFusionClient *client.LightFusionClient, leadRepo *repo.LeadRepo) *Project3DHandler {
	return &Project3DHandler{
		lightFusionClient: lightFusionClient,
		leadRepo:          leadRepo,
	}
}

// Create3DProjectRequest represents the API request for creating a 3D project
// @Description Request body for creating a 3D solar project with energy calculations
type Create3DProjectRequest struct {
	LeadID            *int             `json:"lead_id,omitempty" example:"123"`
	Latitude          float64          `json:"latitude" example:"37.7749"`
	Longitude         float64          `json:"longitude" example:"-122.4194"`
	Address           AddressRequest   `json:"address"`
	Homeowner         HomeownerRequest `json:"homeowner"`
	Hardware          HardwareRequest  `json:"hardware"`
	Consumption       []int            `json:"consumption" example:"800,850,900,950,1000,1050,1100,1150,1200,1250,1300,1350"`
	LseID             int              `json:"lse_id" example:"1"`
	Period            string           `json:"period" example:"month"`
	TargetSolarOffset int              `json:"target_solar_offset" example:"100"`
	Mode              *string          `json:"mode,omitempty" example:"max"`
	Unit              string           `json:"unit" example:"kwh"`
	CompanyID         int              `json:"company_id" example:"1"`
	CreatorID         int              `json:"creator_id" example:"1"`
}

type AddressRequest struct {
	Street     string `json:"street" example:"123 Solar Street"`
	City       string `json:"city" example:"San Francisco"`
	State      string `json:"state" example:"CA"`
	PostalCode string `json:"Postalcode" example:"94102"`
	Country    string `json:"country" example:"USA"`
}

type HomeownerRequest struct {
	Email     string `json:"email" example:"homeowner@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Phone     string `json:"phone" example:"+1234567890"`
}

type HardwareRequest struct {
	PanelID         int  `json:"panel_id" example:"1"`
	InverterID      int  `json:"inverter_id" example:"1"`
	StorageID       *int `json:"storage_id,omitempty" example:"1"`
	StorageQuantity *int `json:"storage_quantity,omitempty" example:"2"`
}

// Create3DProjectResponse represents the API response
type Create3DProjectResponse struct {
	ID               int     `json:"id" example:"123"`
	LeadID           int     `json:"lead_id" example:"456"`
	Status           string  `json:"status" example:"processing"`
	AnnualProduction float64 `json:"annual_production,omitempty" example:"15000"`
	SystemSize       float64 `json:"system_size,omitempty" example:"10.5"`
	EstimatedCost    float64 `json:"estimated_cost,omitempty" example:"25000"`
	AnnualSavings    float64 `json:"annual_savings,omitempty" example:"2500"`
	Message          string  `json:"message" example:"3D project created successfully. Processing in background."`
}

// Create3DProject godoc
// @Summary Create a 3D solar project with energy calculations
// @Description Creates a 3D model from Google Earth data and calculates energy requirements and costs
// @Tags projects
// @Accept json
// @Produce json
// @Param request body Create3DProjectRequest true "Project details"
// @Success 201 {object} Create3DProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/external [post]
func (h *Project3DHandler) Create3DProject(w http.ResponseWriter, r *http.Request) {
	var req Create3DProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		respondError(w, http.StatusBadRequest, "Latitude and longitude are required")
		return
	}

	if req.Address.Street == "" || req.Address.City == "" {
		respondError(w, http.StatusBadRequest, "Address details are required")
		return
	}

	lightFusionReq := client.Create3DProjectRequest{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Address: client.AddressDetails{
			Street:     req.Address.Street,
			City:       req.Address.City,
			State:      req.Address.State,
			PostalCode: req.Address.PostalCode,
			Country:    req.Address.Country,
		},
		Homeowner: client.HomeownerDetails{
			Email:     req.Homeowner.Email,
			FirstName: req.Homeowner.FirstName,
			LastName:  req.Homeowner.LastName,
			Phone:     req.Homeowner.Phone,
		},
		Hardware: client.HardwareDetails{
			PanelID:         req.Hardware.PanelID,
			InverterID:      req.Hardware.InverterID,
			StorageID:       req.Hardware.StorageID,
			StorageQuantity: req.Hardware.StorageQuantity,
		},
		Consumption:       req.Consumption,
		LseID:             req.LseID,
		Period:            req.Period,
		TargetSolarOffset: req.TargetSolarOffset,
		Mode:              req.Mode,
		Unit:              req.Unit,
	}

	resp, err := h.lightFusionClient.Create3DProject(r.Context(), lightFusionReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create 3D project: "+err.Error())
		return
	}
	fmt.Printf("\n\n Response: %v\n", resp)


	var lead *models.Lead
	if req.LeadID != nil {
		lead, err = h.leadRepo.GetByID(r.Context(), *req.LeadID)
		if err != nil {
			log.Printf("Warning: Could not find lead with ID %d to update: %v", *req.LeadID, err)
		} else {
			lead.SetLightFusion3DProject(resp.ID, resp.LeadID)
			if err := h.leadRepo.Update(r.Context(), lead); err != nil {
				log.Printf("Warning: Failed to update lead with 3D project info: %v", err)
			}
		}
	} else {
		lead = &models.Lead{
			CompanyID:        req.CompanyID,
			CreatorID:        req.CreatorID,
			Latitude:         req.Latitude,
			Longitude:        req.Longitude,
			Address:          fmt.Sprintf("%s, %s, %s %s", req.Address.Street, req.Address.City, req.Address.State, req.Address.PostalCode),
			State:            int(models.LeadStateInitialized),
			Source:           int(models.LeadSourceEarth),
			ExternalLeadID:   &resp.LeadID,
			SystemSize:       resp.SystemSize,
			AnnualProduction: resp.AnnualProduction,
		}
		lead.SetLightFusion3DProject(resp.ID, resp.LeadID)

		if err := h.leadRepo.Create(r.Context(), lead); err != nil {
			log.Printf("Warning: Failed to create lead with 3D project info: %v", err)
		}
	}

	response := Create3DProjectResponse{
		ID:               resp.ID,
		LeadID:           resp.LeadID,
		Status:           resp.Status,
		AnnualProduction: resp.AnnualProduction,
		SystemSize:       resp.SystemSize,
		EstimatedCost:    resp.EstimatedCost,
		AnnualSavings:    resp.AnnualSavings,
		Message:          "3D project created successfully. Processing in background.",
	}

	respondJSON(w, http.StatusCreated, response)
}

// GetProjectStatus godoc
// @Summary Get 3D project status
// @Description Retrieves the status and details of a 3D solar project
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Param house_id query int true "House ID"
// @Success 200 {object} client.Status3DProjectResponse "Response structure from LightFusion API"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/external/{id} [get]
func (h *Project3DHandler) GetProjectStatus(w http.ResponseWriter, r *http.Request) {
	var projectID int
	_, err := fmt.Sscanf(r.URL.Path, "/api/projects/3d/%d", &projectID)
	if err != nil {
		_, err = fmt.Sscanf(r.URL.Path, "/projects/3d/%d", &projectID)
		if err != nil {
			errMsg := fmt.Sprintf("Invalid project ID in path '%s': %v", r.URL.Path, err)
			log.Printf("Error: %s", errMsg)
			respondError(w, http.StatusBadRequest, errMsg)
			return
		}
	}

	if projectID == 0 {
		errMsg := "Project ID cannot be 0"
		log.Printf("Error: %s", errMsg)
		respondError(w, http.StatusBadRequest, errMsg)
		return
	}
	houseIDStr := r.URL.Query().Get("house_id")
	houseID, err := strconv.Atoi(houseIDStr)
	if err != nil || houseID == 0 {
		errMsg := fmt.Sprintf("Invalid or missing house_id query param: '%s'", houseIDStr)
		log.Printf("Error: %s", errMsg)
		respondError(w, http.StatusBadRequest, errMsg)
		return
	}
	if houseID == 0 {
		errMsg := "House ID cannot be 0"
		log.Printf("Error: %s", errMsg)
		respondError(w, http.StatusBadRequest, errMsg)
		return
	}

	resp, err := h.lightFusionClient.GetProjectStatus(r.Context(), projectID, houseID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get project status for ID %d: %v", projectID, err)
		log.Printf("Error: %s", errMsg)

		statusCode := http.StatusInternalServerError
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(strings.ToLower(err.Error()), "unauthorized") ||
			strings.Contains(strings.ToLower(err.Error()), "not authenticated") {
			statusCode = http.StatusUnauthorized
		}

		respondError(w, statusCode, errMsg)
		return
	}

	log.Printf("Successfully retrieved project status for ID %d", projectID)

	if resp.LeadCompletion != nil {
		leadData := resp.LeadCompletion.Lead
		lead, err := h.leadRepo.GetByExternalID(r.Context(), leadData.ID)
		if err == nil && lead != nil {
			if leadData.State == 1 {
				lead.Update3DModelStatus("completed")
			} else if leadData.State == 2 {
				lead.Update3DModelStatus("failed")
			} else {
				lead.Update3DModelStatus("processing")
			}
			if leadData.House.SystemSize > 0 {
				lead.SystemSize = float64(leadData.House.SystemSize)
			}
			if leadData.Production.Annual > 0 {
				lead.AnnualProduction = leadData.Production.Annual
			}

			if err := h.leadRepo.Update(r.Context(), lead); err != nil {
				log.Printf("Warning: Failed to update lead with status: %v", err)
			}
		}
	}

	respondJSON(w, http.StatusOK, resp)
}

// GetProjectFiles3D godoc
// @Summary Get 3D project mesh files
// @Description Downloads and retrieves 3D mesh files (JPG, OBJ, PLY, MTL) for a project
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} client.ProfilesFiles3DResponse "3D mesh files response"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/external/{id}/files [get]
func (h *Project3DHandler) GetProjectFiles3D(w http.ResponseWriter, r *http.Request) {
	var projectID int
	_, err := fmt.Sscanf(r.URL.Path, "/api/projects/3d/%d/files", &projectID)
	if err != nil {
		_, err = fmt.Sscanf(r.URL.Path, "/projects/3d/%d/files", &projectID)
		if err != nil {
			errMsg := fmt.Sprintf("Invalid project ID in path '%s': %v", r.URL.Path, err)
			log.Printf("Error: %s", errMsg)
			respondError(w, http.StatusBadRequest, errMsg)
			return
		}
	}
	if projectID == 0 {
		errMsg := "Project ID cannot be 0"
		log.Printf("Error: %s", errMsg)
		respondError(w, http.StatusBadRequest, errMsg)
		return
	}

	resp, err := h.lightFusionClient.GetProjectFiles(r.Context(), projectID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get project files for ID %d: %v", projectID, err)
		log.Printf("Error: %s", errMsg)

		statusCode := http.StatusInternalServerError
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(strings.ToLower(err.Error()), "unauthorized") ||
			strings.Contains(strings.ToLower(err.Error()), "not authenticated") {
			statusCode = http.StatusUnauthorized
		}

		respondError(w, statusCode, errMsg)
		return
	}

	log.Printf("Successfully retrieved project files for ID %d", projectID)

	respondJSON(w, http.StatusOK, resp)
}
