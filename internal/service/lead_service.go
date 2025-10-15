package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
)

type CreateLeadResponse struct {
	Success bool `json:"success"`
	LeadID  int  `json:"lead_id"`
	HouseID int  `json:"house_id"`
}

type LeadService struct {
	leadRepo *repo.LeadRepo
	houseRepo *repo.HouseRepo
	genabilityClient *client.Agent
}

type CreateLead struct {
	CompanyID        int     `json:"company_id" example:"1"`
	CreatorID        int     `json:"creator_id" example:"1"`
	Latitude         float64 `json:"latitude" example:"37.7749"`
	Longitude        float64 `json:"longitude" example:"-122.4194"`
	Address          string  `json:"address" example:"123 Solar St, San Francisco, CA 94102"`
	Street           *string `json:"street,omitempty"`
	City             *string `json:"city,omitempty"`
	State            *string `json:"state,omitempty"`
	Zip              *string `json:"zip,omitempty"`
	HomeownerName    *string `json:"homeowner_name,omitempty"`
	HomeownerEmail   *string `json:"homeowner_email,omitempty"`
	HomeownerPhone   *string `json:"homeowner_phone,omitempty"`
	SystemSize       float64 `json:"system_size" example:"10.5"`
	PanelCount       int     `json:"panel_count" example:"30"`
	HardwareType     *string `json:"hardware_type,omitempty"`
	KwhUsage         float64 `json:"kwh_usage" example:"12000"`
	Consumption      []int   `json:"consumption,omitempty"`
	SalesRepEmail     *string `json:"sales_rep_email,omitempty"`
	LseId             int     `json:"lse_id"`
	Period            string  `json:"period"`
	TargetSolarOffset int     `json:"target_solar_offset"`
	Mode              *string `json:"mode,omitempty"`
	Unit              string  `json:"unit"`
}


func NewLeadService(leadRepo *repo.LeadRepo, houseRepo *repo.HouseRepo) *LeadService {
	var genClient *client.Agent

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed to initialize Genability client: %v", r)
		}
	}()

	genClient = client.NewAgent()

	return &LeadService{
		leadRepo: leadRepo,
		houseRepo: houseRepo,
		genabilityClient: genClient,
	}
}

func (s *LeadService) CreateLead(ctx context.Context, req CreateLead, userID int, effectiveCompanyID int) (*CreateLeadResponse, error) {
	house := models.House{
			Lat:         req.Latitude,
			Lng:         req.Longitude,
			Diameter:    20,
			Probability: 0.70,
			State:       "pending",
		}
	if err := s.houseRepo.Create(ctx, &house); err != nil {
		return nil, err
	}
	houseID := house.ID

	lead := models.Lead{
			CompanyID:  effectiveCompanyID,
			CreatorID:  userID,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			Address:    req.Address,
			KwhUsage:   req.KwhUsage,
			SystemSize: req.SystemSize,
			PanelCount: req.PanelCount,
			Source:     0,
			State:      0,
		}
		if err := s.leadRepo.Create(ctx, &lead); err != nil {
			return nil, fmt.Errorf("failed to create lead: %w", err)
		}

		if s.genabilityClient != nil {
			accountInput := client.Account{
				Address: client.AccountAddress{
					String:    lead.Address,
					Latitude:  lead.Latitude,
					Longitude: lead.Longitude,
				},
			}
			accounts := client.NewAccounts(s.genabilityClient)
			genAcc, err := accounts.Create(ctx, accountInput)
			if err != nil {
				log.Printf("Warning: Failed to create Genability account: %v", err)
			} else if genAcc != nil {
				tariffs := client.NewTariffs(s.genabilityClient)
				tariff, err := tariffs.GetCurrent(ctx, genAcc.ID)
				if err != nil {
					log.Printf("Warning: Failed to get current tariff: %v", err)
				} else if tariff != nil {
					utilityID := int(tariff.LseID)
					tariffID := int(tariff.ID)
					lead.UtilityID = &utilityID
					lead.TariffID = &tariffID

					if err := s.leadRepo.Update(ctx, &lead); err != nil {
						log.Printf("Warning: Failed to update lead with Genability data: %v", err)
					}
				}
			}
		}

	return &CreateLeadResponse{
		Success: true,
		LeadID:  lead.ID,
		HouseID: int(houseID),
	}, nil
}
