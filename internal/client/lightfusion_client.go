package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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
		baseURL:    baseURL,
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

// LoginRequest represents the login request for LightFusion API
type LoginRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
	Expires  bool   `json:"expires"`
}

// LoginResponse represents the login response from LightFusion API
type LoginResponse struct {
	Token   string      `json:"token"`
	User    interface{} `json:"user"`
	Contact interface{} `json:"contact"`
	Company interface{} `json:"company"`
}

// Login authenticates with the LightFusion API and returns a session token
func (c *LightFusionClient) Login(ctx context.Context, email, password string) (string, error) {
	// Create a simple form data request
	formData := fmt.Sprintf(`{"contact":"%s","password":"%s"}`, email, password)
	body := bytes.NewBufferString(formData)

	// Log the request
	endpoint := c.baseURL + "/v1/users/sessions"
	log.Printf("Sending login request to %s with body: %s", endpoint, formData)

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Log the full request
	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("Request details:\n%s", string(dump))

	// Create a new HTTP client with disabled SSL verification (for testing only)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(respBody))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var result struct {
		Token   string      `json:"token"`
		User    interface{} `json:"user"`
		Contact interface{} `json:"contact"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Token == "" {
		return "", fmt.Errorf("no token in response: %s", string(respBody))
	}

	// Store the token for future requests
	c.apiKey = result.Token

	return result.Token, nil
}

// Create3DProjectRequest represents the request to create a 3D project with energy calculations
type Create3DProjectRequest struct {
	Latitude          float64          `json:"latitude"`
	Longitude         float64          `json:"longitude"`
	Address           AddressDetails   `json:"address"`
	Homeowner         HomeownerDetails `json:"homeowner"`
	Hardware          HardwareDetails  `json:"hardware"`
	Consumption       []int            `json:"consumption"`
	LseID             int              `json:"lseId"`
	Period            string           `json:"period"`
	TargetSolarOffset int              `json:"targetSolarOffset"`
	Mode              *string          `json:"mode,omitempty"`
	Unit              string           `json:"unit"`
}

// AddressDetails represents address information
type AddressDetails struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// HomeownerDetails represents homeowner information
type HomeownerDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
}

// HardwareDetails represents solar hardware configuration
type HardwareDetails struct {
	PanelID         int  `json:"panel_id"`
	InverterID      int  `json:"inverter_id"`
	StorageID       *int `json:"storage_id,omitempty"`
	StorageQuantity *int `json:"storage_quantity,omitempty"`
}

// Create3DProjectResponse represents the response from 3D project creation
type Create3DProjectResponse struct {
	ID               int     `json:"id"`
	LeadID           int     `json:"lead_id"`
	Status           string  `json:"status"`
	AnnualProduction float64 `json:"annual_production,omitempty"`
	SystemSize       float64 `json:"system_size,omitempty"`
	EstimatedCost    float64 `json:"estimated_cost,omitempty"`
	AnnualSavings    float64 `json:"annual_savings,omitempty"`
}

// Create3DProject creates a 3D project with energy calculations using the LightFusion API
func (c *LightFusionClient) Create3DProject(ctx context.Context, req Create3DProjectRequest) (*Create3DProjectResponse, error) {

	if c.apiKey == "" {
		return nil, fmt.Errorf("not authenticated with LightFusion API")
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("Creating 3D project with request: %s", string(reqJSON))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/lead/create", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	dump, _ := httputil.DumpRequestOut(httpReq, true)
	log.Printf("Sending request to create 3D project:\n%s", string(dump))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %d, body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var projectResp Create3DProjectResponse
	if err := json.Unmarshal(bodyBytes, &projectResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	return &projectResp, nil
}

// GetProjectStatus retrieves the status of a 3D project
func (c *LightFusionClient) GetProjectStatus(ctx context.Context, projectID int) (*Create3DProjectResponse, error) {

	if c.apiKey == "" {
		return nil, fmt.Errorf("not authenticated with LightFusion API")
	}

	endpoint := fmt.Sprintf("%s/v1/projects/3d/%d", c.baseURL, projectID)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	log.Printf("Fetching project status from %s", endpoint)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Project status response status: %d, body: %s", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var projectResp Create3DProjectResponse
	if err := json.Unmarshal(bodyBytes, &projectResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	return &projectResp, nil
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
	ID               int       `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	State            int       `json:"state"`
	CompanyID        int       `json:"company_id"`
	CreatorID        int       `json:"creator_id"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	Address          string    `json:"address"`
	Source           int       `json:"source"`
	PromoCode        *string   `json:"promo_code"`
	Is2D             bool      `json:"is_2d"`
	KwhUsage         float64   `json:"kwh_usage"`
	SystemSize       float64   `json:"system_size"`
	PanelCount       int       `json:"panel_count"`
	PanelID          *int      `json:"panel_id"`
	InverterID       *int      `json:"inverter_id"`
	UtilityID        *int      `json:"utility_id"`
	RoofMaterial     *int      `json:"roof_material"`
	AnnualProduction float64   `json:"annual_production"`
	InstallationDate *string   `json:"installation_date"`
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
