package client


import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	AccountPropertyCustomerClass  = "customerClass"
	AccountPropertyLseID          = "lseId"
	AccountPropertyBuildingArea   = "buildingArea"
	AccountPropertyMasterTariffID = "masterTariffId"
	AccountPropertyTerritoryID    = "territoryId"
)

type Credentials struct {
	AppID  string
	AppKey string
}

type Agent struct {
	client *http.Client
	creds  Credentials
	base   string
}

type AccountAddress struct {
	String     string  `json:"addressString,omitempty"`
	Street     string  `json:"address1,omitempty"`
	City       string  `json:"city,omitempty"`
	State      string  `json:"state,omitempty"`
	Country    string  `json:"country,omitempty"`
	Postalcode string  `json:"zip,omitempty"`
	Latitude   float64 `json:"lat,omitempty"`
	Longitude  float64 `json:"lon,omitempty"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AccountProperty struct {
	KeyValue `json:",inline"`
	Type     string  `json:"dataType,omitempty"`
	Accuracy float64 `json:"accuracy,omitempty"`
	Name     string  `json:"displayName,omitempty"`
}

type Tariff struct {
	ID          uint    `json:"tariffId"`
	Name        string  `json:"tariffName"`
	Code        string  `json:"tariffCode"`
	LseID       uint    `json:"lseId"`
	LseName     string  `json:"lseName"`
	ServiceType string  `json:"serviceType"`
	IsActive    bool    `json:"isActive"`
	MasterID    uint    `json:"masterTariffId"`
	CustomerLikelihood *float64 `json:"customerLikelihood,omitempty"`
}

type Tariffs struct {
	Agent *Agent
	cache *cache.Cache
}



func NewAgent() *Agent {
	appID := os.Getenv("GENABILITY_ID")
	appKey := os.Getenv("GENABILITY_KEY")

	if appID == "" || appKey == "" {
		log.Fatal("missing GENABILITY_ID or GENABILITY_KEY in environment")
	}
	Creds := Credentials{AppID: appID, AppKey: appKey}
	return &Agent{
		client: &http.Client{Timeout: 10 * time.Second},
		creds:  Creds,
		base:   "https://api.genability.com/rest/",
	}
}

func (a *Agent) doRequest(ctx context.Context, method, path string, body any, out any) error {
	url := a.base + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(a.creds.AppID, a.creds.AppKey)
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("genability: request failed with status %d", resp.StatusCode)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}



func NewAccountProp(key, val string) AccountProperty {
	return AccountProperty{
		KeyValue: KeyValue{Key: key, Value: val},
	}
}

type Account struct {
	ID         string                     `json:"accountId"`
	Name       string                     `json:"accountName"`
	Status     string                     `json:"status"`
	Address    AccountAddress             `json:"address"`
	Properties map[string]AccountProperty `json:"properties"`
}

type Accounts struct {
	Agent *Agent
}

func NewAccounts(agent *Agent) *Accounts {
	return &Accounts{Agent: agent}
}

func (a *Accounts) Create(ctx context.Context, input Account) (*Account, error) {
	var resp struct {
		Results []Account `json:"results"`
	}
	if err := a.Agent.doRequest(ctx, "POST", "v1/accounts", input, &resp); err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("no account returned from Genability")
	}
	return &resp.Results[0], nil
}

func (a *Accounts) Show(ctx context.Context, id string) (*Account, error) {
	var resp struct {
		Results []Account `json:"results"`
	}
	if err := a.Agent.doRequest(ctx, "GET", "v1/accounts/"+id, nil, &resp); err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("account not found")
	}
	return &resp.Results[0], nil
}


func NewTariffs(agent *Agent) *Tariffs {
	return &Tariffs{
		Agent: agent,
		cache: cache.New(time.Hour, cache.NoExpiration),
	}
}

func (t *Tariffs) fetch(ctx context.Context, v url.Values) ([]Tariff, error) {
	page := 0
	v.Set("isActive", "true")
	v.Set("pageCount", "100")
	v.Set("customerClasses", "RESIDENTIAL")

	var all []Tariff
	
	for {
		v.Set("pageStart", strconv.Itoa(page*100))
		page++

		var resp struct {
			Results []Tariff `json:"results"`
			Count   int      `json:"count"`
		}
		if err := t.Agent.doRequest(ctx, "GET", "public/tariffs?"+v.Encode(), nil, &resp); err != nil {
			return nil, err
		}

		all = append(all, resp.Results...)
		if page*100 >= resp.Count {
			break
		}
	}
	return all, nil
}

func (t *Tariffs) Index(ctx context.Context, zipcode, country string) ([]Tariff, error) {
	key := fmt.Sprintf("zip_%s_%s", zipcode, country)
	if cached, ok := t.cache.Get(key); ok {
		return cached.([]Tariff), nil
	}
	v := url.Values{}
	v.Set("zipCode", zipcode)
	v.Set("country", country)

	data, err := t.fetch(ctx, v)
	if err == nil {
		t.cache.Set(key, data, cache.DefaultExpiration)
	}
	return data, err
}

func (t *Tariffs) Show(ctx context.Context, masterID uint) (*Tariff, error) {
	key := fmt.Sprintf("tariff_%d", masterID)
	if cached, ok := t.cache.Get(key); ok {
		tariff := cached.(Tariff)
		return &tariff, nil
	}
	
	var resp struct {
		Results []Tariff `json:"results"`
	}
	if err := t.Agent.doRequest(ctx, "GET", "public/tariffs/"+strconv.Itoa(int(masterID)), nil, &resp); err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("tariff not found")
	}
	tariff := resp.Results[0]
	t.cache.Set(key, tariff, cache.DefaultExpiration)
	return &tariff, nil
}

// GetCurrent retrieves the current tariff for a given account ID
func (t *Tariffs) GetCurrent(ctx context.Context, accountID string) (*Tariff, error) {
	key := fmt.Sprintf("account_tariff_%s", accountID)
	if cached, ok := t.cache.Get(key); ok {
		tariff := cached.(Tariff)
		return &tariff, nil
	}
	
	var resp struct {
		Results []Tariff `json:"results"`
		Count   int      `json:"count"`
	}
	
	path := fmt.Sprintf("v1/accounts/%s/tariffs", accountID)
	if err := t.Agent.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	
	if len(resp.Results) == 0 {
		return nil, errors.New("no tariff found for account")
	}
	
	// Return the first active tariff, or the first one if none are active
	for _, tariff := range resp.Results {
		if tariff.IsActive {
			t.cache.Set(key, tariff, cache.DefaultExpiration)
			return &tariff, nil
		}
	}
	
	// If no active tariff found, return the first one
	tariff := resp.Results[0]
	t.cache.Set(key, tariff, cache.DefaultExpiration)
	return &tariff, nil
}
