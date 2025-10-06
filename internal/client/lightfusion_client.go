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
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LightFusionClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewLightFusionClient(baseURL, apiKey string) *LightFusionClient {
	return &LightFusionClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

type LoginRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
	Expires  bool   `json:"expires"`
}

type LoginResponse struct {
	Token   string      `json:"token"`
	User    interface{} `json:"user"`
	Contact interface{} `json:"contact"`
	Company interface{} `json:"company"`
}

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

type AddressDetails struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type HomeownerDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type HardwareDetails struct {
	PanelID         int  `json:"panel_id"`
	InverterID      int  `json:"inverter_id"`
	StorageID       *int `json:"storage_id,omitempty"`
	StorageQuantity *int `json:"storage_quantity,omitempty"`
}

type Create3DProjectResponse struct {
	ID               int                    `json:"id"`
	LeadID           int                    `json:"lead_id"`
	Status           string                 `json:"status"`
	AnnualProduction float64                `json:"annual_production,omitempty"`
	SystemSize       float64                `json:"system_size,omitempty"`
	EstimatedCost    float64                `json:"estimated_cost,omitempty"`
	AnnualSavings    float64                `json:"annual_savings,omitempty"`
	Adders           map[string]interface{} `json:"adders,omitempty"`
}



type Status3DProjectResponse struct {
    Panel    *Panel    `json:"panel"`
    Inverter []Inverter `json:"inverter"`
    Adders   []Adder   `json:"adders"`
    PriceBreakdown *PriceBreakdown `json:"price_breakdown,omitempty"`
}

type Panel struct {
    ID           int     `json:"ID"`
    Manufacturer string  `json:"Manufacturer"`
    Model        string  `json:"Model"`
    DisplayName  string  `json:"DisplayName"`
    Active       bool    `json:"Active"`
    IsDefault    bool    `json:"IsDefault"`
    Power        float64 `json:"Power"`
    PricePerWatt float64 `json:"PricePerWatt"`
    IsDomestic   bool    `json:"IsDomestic"`
}

type Inverter struct {
    Name         string  `json:"Name"`
    CostType     string  `json:"CostType"`
    Category     string  `json:"Category"`
    IsActive     bool    `json:"IsActive"`
    Cost         float64 `json:"Cost"`
    ID           int     `json:"ID"`
    Manufacturer string  `json:"Manufacturer"`
    Capacity     float64 `json:"Capacity"`
    Quantity     int     `json:"Quantity"`
}


type Adder struct {
    ID          int      `json:"ID"`
    CompanyID   int      `json:"CompanyID"`
    Name        string   `json:"Name"`
    Cost        float64  `json:"Cost"`
    CostType    string   `json:"CostType"`
    States      []string `json:"States"`
    Active      bool     `json:"Active"`
    CategoryID  int      `json:"CategoryID"`
    IsAutomatic bool     `json:"IsAutomatic"`
    MinSystemSize float64 `json:"MinSystemSize"`
    MaxSystemSize float64 `json:"MaxSystemSize"`
    Quantity    int      `json:"Quantity"`
    CustomPrice float64  `json:"CustomPrice"`
}

type PriceBreakdown struct {
	Items                       []PriceItem `json:"items"`
	BasePricePerWatt            float64     `json:"base_price_per_watt"`
	TotalPricePerWatt           float64     `json:"total_price_per_watt"`
	TotalPricePerWattFinanced   float64     `json:"total_price_per_watt_financed"`
	DefaultBasePrice            float64     `json:"default_base_price"`
	MinimumBasePrice            float64     `json:"minimum_base_price"`
	TotalAmount                 float64     `json:"total_amount"`
	TotalAmountWithoutDealerFee float64     `json:"total_amount_without_dealer_fee"`
	TotalFee                    float64     `json:"total_fee"`
}

type PriceItem struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProfilesFiles3DResponse struct {
	ProjectID  int      `json:"project_id"`
	JPGPath    string   `json:"jpg_path"`
	OBJPath    string   `json:"obj_path"`
	PLYPath    string   `json:"ply_path"`
	MTLPath    string   `json:"mtl_path"`
	JPGURL     string   `json:"jpg_url"`
	OBJURL     string   `json:"obj_url"`
	PLYURL     string   `json:"ply_url"`
	MTLURL     string   `json:"mtl_url"`
	Downloaded bool     `json:"downloaded"`
	Errors     []string `json:"errors,omitempty"`
}

func (c *LightFusionClient) Login(ctx context.Context, email, password string) (string, error) {
	formData := fmt.Sprintf(`{"contact":"%s","password":"%s"}`, email, password)
	body := bytes.NewBufferString(formData)

	endpoint := c.baseURL + "/v1/users/sessions"
	log.Printf("Sending login request to %s with body: %s", endpoint, formData)

	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(respBody))
	}

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

	c.apiKey = result.Token

	return result.Token, nil
}

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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var projectResp Create3DProjectResponse
	if err := json.Unmarshal(bodyBytes, &projectResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	return &projectResp, nil
}

func (c *LightFusionClient) GetProjectStatus(ctx context.Context, projectID int, houseID int) (*Status3DProjectResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("not authenticated with LightFusion API")
	}

	requestBody := struct {
		ProjectID int `json:"project_id"`
	}{
		ProjectID: projectID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v3/adders.ListProjectAdders", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	log.Printf("Fetching project adders from %s for project ID: %d", endpoint, projectID)

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


	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var projectResp Status3DProjectResponse
	if err := json.Unmarshal(bodyBytes, &projectResp); err != nil {
		return nil, fmt.Errorf(" hookah bar failed to decode response: %w, body: %s", err, string(bodyBytes))
	}

	priceRequestBody := struct {
		ProjectID int `json:"project_id"`
		HouseID   int `json:"house_id"`
	}{
		ProjectID: projectID,
		HouseID:   houseID,
	}

	priceJsonBody, err := json.Marshal(priceRequestBody)
	if err != nil {
		log.Printf("Warning: failed to marshal price breakdown request: %v", err)
	} else {
		priceEndpoint := fmt.Sprintf("%s/v3/adders.GetPriceBreakdown", c.baseURL)
		priceReq, err := http.NewRequestWithContext(ctx, "POST", priceEndpoint, bytes.NewBuffer(priceJsonBody))
		if err != nil {
			log.Printf("Warning: failed to create price breakdown request: %v", err)
		} else {
			priceReq.Header.Set("Content-Type", "application/json")
			priceReq.Header.Set("Accept", "application/json")
			priceReq.Header.Set("Authorization", "Bearer "+c.apiKey)

			log.Printf("Fetching price breakdown from %s for project ID: %d, house ID: %d", priceEndpoint, projectID, houseID)

			priceResp, err := client.Do(priceReq)
			if err != nil {
				log.Printf("Warning: failed to fetch price breakdown: %v", err)
			} else {
				defer priceResp.Body.Close()
				priceBodyBytes, _ := io.ReadAll(priceResp.Body)
				if priceResp.StatusCode == http.StatusOK {
					var priceBreakdown PriceBreakdown
					if err := json.Unmarshal(priceBodyBytes, &priceBreakdown); err != nil {
						log.Printf("Warning: failed to decode price breakdown: %v", err)
					} else {
						projectResp.PriceBreakdown = &priceBreakdown
					}
				} else {
					log.Printf("Warning: price breakdown API returned status %d: %s", priceResp.StatusCode, string(priceBodyBytes))
				}
			}
		}
	}

	return &projectResp, nil
}

func (c *LightFusionClient) GetProjectFiles(ctx context.Context, projectID int) (*ProfilesFiles3DResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("not authenticated with LightFusion API")
	}

	response := &ProfilesFiles3DResponse{
		ProjectID: projectID,
		Errors:    []string{},
	}

	mediaDir := "./media"
	projectDir := filepath.Join(mediaDir, fmt.Sprintf("%d", projectID))
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	files := []struct {
		filename string
		path     *string
		url      *string
	}{
		{"scene.jpg", &response.JPGPath, &response.JPGURL},
		{"scene.obj", &response.OBJPath, &response.OBJURL},
		{"scene.ply", &response.PLYPath, &response.PLYURL},
		{"scene.mtl", &response.MTLPath, &response.MTLURL},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0

	for _, file := range files {
		wg.Add(1)
		go func(f struct {
			filename string
			path     *string
			url      *string
		}) {
			defer wg.Done()

			filePath := filepath.Join(projectDir, f.filename)
			fileURL := fmt.Sprintf("/media/%d/%s", projectID, f.filename)

			if _, err := os.Stat(filePath); err == nil {
				log.Printf("File already exists: %s", filePath)
				mu.Lock()
				*f.path = filePath
				*f.url = fileURL
				successCount++
				mu.Unlock()
				return
			}

			endpoint := fmt.Sprintf("https://storage.googleapis.com/lightfusiondev/leads/%d/mesh/%s", projectID, f.filename)
			if err := c.downloadMeshFile(ctx, endpoint, filePath); err != nil {
				errMsg := fmt.Sprintf("failed to download %s: %v", f.filename, err)
				log.Printf("Warning: %s", errMsg)
				mu.Lock()
				response.Errors = append(response.Errors, errMsg)
				mu.Unlock()
				return
			}

			mu.Lock()
			*f.path = filePath
			*f.url = fileURL
			successCount++
			mu.Unlock()

			log.Printf("Successfully downloaded %s to %s", f.filename, filePath)
		}(file)
	}

	wg.Wait()

	response.Downloaded = successCount > 0

	if successCount == 0 {
		return response, fmt.Errorf("failed to download any mesh files")
	}

	return response, nil
}

func (c *LightFusionClient) downloadMeshFile(ctx context.Context, endpoint, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	if strings.Contains(endpoint, "api.lightfusion.io") {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	log.Printf("Downloading file from %s", endpoint)
	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("Outgoing request:\n%s", string(dump))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Downloaded %d bytes to %s", written, destPath)

	return nil
}
