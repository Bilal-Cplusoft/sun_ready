package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
)

type Project3DHandler struct {
	lightFusionClient *client.LightFusionClient
}

func NewProject3DHandler(lightFusionClient *client.LightFusionClient) *Project3DHandler {
	return &Project3DHandler{
		lightFusionClient: lightFusionClient,
	}
}

// Create3DProjectRequest represents the API request for creating a 3D project
// @Description Request body for creating a 3D solar project with energy calculations
type Create3DProjectRequest struct {
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
// @Router /projects/3d [post]
func (h *Project3DHandler) Create3DProject(w http.ResponseWriter, r *http.Request) {
	var req Create3DProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Latitude == 0 || req.Longitude == 0 {
		respondWithError(w, http.StatusBadRequest, "Latitude and longitude are required")
		return
	}

	if req.Address.Street == "" || req.Address.City == "" {
		respondWithError(w, http.StatusBadRequest, "Address details are required")
		return
	}

	// Convert to LightFusion API format
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
		respondWithError(w, http.StatusInternalServerError, "Failed to create 3D project: "+err.Error())
		return
	}
	fmt.Printf("\n\n Response: %v\n", resp)

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

	respondWithJSON(w, http.StatusCreated, response)
}

// GetProjectStatus godoc
// @Summary Get 3D project status
// @Description Retrieves the status and details of a 3D solar project
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} Create3DProjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/projects/3d/{id} [get]
func (h *Project3DHandler) GetProjectStatus(w http.ResponseWriter, r *http.Request) {
	// Extract project ID from URL path
	var projectID int
	_, err := fmt.Sscanf(r.URL.Path, "/api/projects/3d/%d", &projectID)
	if err != nil {
		// Try alternative path format
		_, err = fmt.Sscanf(r.URL.Path, "/projects/3d/%d", &projectID)
		if err != nil {
			errMsg := fmt.Sprintf("Invalid project ID in path '%s': %v", r.URL.Path, err)
			log.Printf("Error: %s", errMsg)
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
	}

	if projectID == 0 {
		errMsg := "Project ID cannot be 0"
		log.Printf("Error: %s", errMsg)
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	log.Printf("Fetching status for project ID: %d", projectID)

	// Call LightFusion API - no need for auth header as client is already authenticated
	resp, err := h.lightFusionClient.GetProjectStatus(r.Context(), projectID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get project status for ID %d: %v", projectID, err)
		log.Printf("Error: %s", errMsg)

		// Return appropriate status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "not authenticated with LightFusion API" {
			statusCode = http.StatusInternalServerError // This should not happen as we're already authenticated
		} else if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}

		respondWithError(w, statusCode, errMsg)
		return
	}

	log.Printf("Successfully retrieved project status for ID %d", projectID)

	// Map the response to our API response format
	response := Create3DProjectResponse{
		ID:               resp.ID,
		LeadID:           resp.LeadID,
		Status:           resp.Status,
		AnnualProduction: resp.AnnualProduction,
		SystemSize:       resp.SystemSize,
		EstimatedCost:    resp.EstimatedCost,
		AnnualSavings:    resp.AnnualSavings,
		Message:          "", // Will be set based on status if needed
	}

	// Add a friendly message based on status
	switch resp.Status {
	case "processing":
		response.Message = "3D project is being processed. Please check back later."
	case "completed":
		response.Message = "3D project processing is complete."
	case "failed":
		response.Message = "3D project processing failed. Please try again or contact support."
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
