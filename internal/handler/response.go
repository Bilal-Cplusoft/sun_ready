package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// LeadResponse represents a successful lead creation response
type LeadResponse struct {
	Success bool `json:"success" example:"true"`
	LeadID  int  `json:"lead_id" example:"42"`
	HouseID int  `json:"house_id" example:"123"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
