package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LightFusionClient handles communication with the LightFUSION API
type LightFusionClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewLightFusionClient creates a new LightFUSION API client
func NewLightFusionClient(baseURL, apiKey string) *LightFusionClient {
	return &LightFusionClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// CreateLeadRequest represents the request to create a lead in LightFUSION
type CreateLeadRequest struct {
	CompanyID    int     `json:"company_id"`
	CreatorID    int     `json:"creator_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Address      string  `json:"address"`
	Source       int     `json:"source"`
	PromoCode    *string `json:"promo_code,omitempty"`
	Is2D         bool    `json:"is_2d"`
	KwhUsage     float64 `json:"kwh_usage"`
	SystemSize   float64 `json:"system_size,omitempty"`
	PanelCount   int     `json:"panel_count,omitempty"`
	PanelID      *int    `json:"panel_id,omitempty"`
	InverterID   *int    `json:"inverter_id,omitempty"`
	UtilityID    *int    `json:"utility_id,omitempty"`
	RoofMaterial *int    `json:"roof_material,omitempty"`
}

// LeadResponse represents a lead from LightFUSION API
type LeadResponse struct {
	ID                  int        `json:"id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	State               int        `json:"state"`
	CompanyID           int        `json:"company_id"`
	CreatorID           int        `json:"creator_id"`
	Latitude            float64    `json:"latitude"`
	Longitude           float64    `json:"longitude"`
	Address             string     `json:"address"`
	Source              int        `json:"source"`
	PromoCode           *string    `json:"promo_code"`
	Is2D                bool       `json:"is_2d"`
	KwhUsage            float64    `json:"kwh_usage"`
	SystemSize          float64    `json:"system_size"`
	PanelCount          int        `json:"panel_count"`
	PanelID             *int       `json:"panel_id"`
	InverterID          *int       `json:"inverter_id"`
	UtilityID           *int       `json:"utility_id"`
	RoofMaterial        *int       `json:"roof_material"`
	AnnualProduction    float64    `json:"annual_production"`
	InstallationDate    *string    `json:"installation_date"`
}

// CreateLead creates a lead in the external LightFUSION API
func (c *LightFusionClient) CreateLead(ctx context.Context, req CreateLeadRequest) (*LeadResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/lead/create", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var leadResp LeadResponse
	if err := json.NewDecoder(resp.Body).Decode(&leadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &leadResp, nil
}

// GetLead retrieves a lead from the external API
func (c *LightFusionClient) GetLead(ctx context.Context, leadID int) (*LeadResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/lead/%d", c.baseURL, leadID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var leadResp LeadResponse
	if err := json.NewDecoder(resp.Body).Decode(&leadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &leadResp, nil
}

// UpdateLead updates a lead in the external API
func (c *LightFusionClient) UpdateLead(ctx context.Context, leadID int, updates map[string]interface{}) (*LeadResponse, error) {
	body, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/v1/lead/%d", c.baseURL, leadID), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var leadResp LeadResponse
	if err := json.NewDecoder(resp.Body).Decode(&leadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &leadResp, nil
}

// ListLeads retrieves leads from the external API
func (c *LightFusionClient) ListLeads(ctx context.Context, companyID int, limit, offset int) ([]LeadResponse, error) {
	url := fmt.Sprintf("%s/v1/leads?company_id=%d&limit=%d&offset=%d", c.baseURL, companyID, limit, offset)
	
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var leadsResp struct {
		Leads []LeadResponse `json:"leads"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&leadsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return leadsResp.Leads, nil
}
