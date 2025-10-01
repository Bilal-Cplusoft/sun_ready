package service

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
)

type LeadService struct {
	leadRepo          *repo.LeadRepo
	lightFusionClient *client.LightFusionClient
	useExternalAPI    bool
}

func NewLeadService(leadRepo *repo.LeadRepo, lightFusionClient *client.LightFusionClient, useExternalAPI bool) *LeadService {
	return &LeadService{
		leadRepo:          leadRepo,
		lightFusionClient: lightFusionClient,
		useExternalAPI:    useExternalAPI,
	}
}

func (s *LeadService) Create(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return err
	}

	if s.useExternalAPI && s.lightFusionClient != nil {
		externalReq := client.CreateLeadRequest{
			CompanyID:    lead.CompanyID,
			CreatorID:    lead.CreatorID,
			Latitude:     lead.Latitude,
			Longitude:    lead.Longitude,
			Address:      lead.Address,
			Source:       lead.Source,
			PromoCode:    lead.PromoCode,
			Is2D:         lead.Is2D,
			KwhUsage:     lead.KwhUsage,
			SystemSize:   lead.SystemSize,
			PanelCount:   lead.PanelCount,
			PanelID:      lead.PanelID,
			InverterID:   lead.InverterID,
			UtilityID:    lead.UtilityID,
			RoofMaterial: lead.RoofMaterial,
		}

		externalLead, err := s.lightFusionClient.CreateLead(ctx, externalReq)
		if err != nil {
			lead.MarkSyncFailed()
		} else {
			lead.MarkSynced(externalLead.ID)
		}
	}

	return s.leadRepo.Create(ctx, lead)
}

func (s *LeadService) GetByID(ctx context.Context, id int) (*models.Lead, error) {
	lead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.useExternalAPI && s.lightFusionClient != nil && lead.ExternalLeadID != nil {
		externalLead, err := s.lightFusionClient.GetLead(ctx, *lead.ExternalLeadID)
		if err == nil {
			s.syncFromExternal(lead, externalLead)
			s.leadRepo.Update(ctx, lead)
		}
	}

	return lead, nil
}

func (s *LeadService) Update(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return err
	}

	if s.useExternalAPI && s.lightFusionClient != nil && lead.ExternalLeadID != nil {
		updates := map[string]interface{}{
			"state":             lead.State,
			"latitude":          lead.Latitude,
			"longitude":         lead.Longitude,
			"address":           lead.Address,
			"kwh_usage":         lead.KwhUsage,
			"system_size":       lead.SystemSize,
			"panel_count":       lead.PanelCount,
			"annual_production": lead.AnnualProduction,
			"installation_date": lead.InstallationDate,
		}

		_, err := s.lightFusionClient.UpdateLead(ctx, *lead.ExternalLeadID, updates)
		if err != nil {
			lead.MarkSyncFailed()
		} else {
			lead.MarkSynced(*lead.ExternalLeadID)
		}
	}

	return s.leadRepo.Update(ctx, lead)
}

func (s *LeadService) Delete(ctx context.Context, id int) error {
	return s.leadRepo.Delete(ctx, id)
}

func (s *LeadService) List(ctx context.Context, limit, offset int) ([]*models.Lead, error) {
	return s.leadRepo.List(ctx, limit, offset)
}

func (s *LeadService) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Lead, error) {
	leads, err := s.leadRepo.ListByCompany(ctx, companyID, limit, offset)
	if err != nil {
		return nil, err
	}

	if s.useExternalAPI && s.lightFusionClient != nil {
		externalLeads, err := s.lightFusionClient.ListLeads(ctx, companyID, limit, offset)
		if err == nil {
			for _, extLead := range externalLeads {
				s.syncExternalLeadToLocal(ctx, &extLead)
			}
			leads, _ = s.leadRepo.ListByCompany(ctx, companyID, limit, offset)
		}
	}

	return leads, nil
}

func (s *LeadService) ListByCreator(ctx context.Context, creatorID int, limit, offset int) ([]*models.Lead, error) {
	return s.leadRepo.ListByCreator(ctx, creatorID, limit, offset)
}

func (s *LeadService) ListByState(ctx context.Context, state int, limit, offset int) ([]*models.Lead, error) {
	return s.leadRepo.ListByState(ctx, state, limit, offset)
}

func (s *LeadService) UpdateState(ctx context.Context, id int, state int) error {
	lead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	lead.State = state

	if s.useExternalAPI && s.lightFusionClient != nil && lead.ExternalLeadID != nil {
		updates := map[string]interface{}{"state": state}
		_, err := s.lightFusionClient.UpdateLead(ctx, *lead.ExternalLeadID, updates)
		if err != nil {
			lead.MarkSyncFailed()
		} else {
			lead.MarkSynced(*lead.ExternalLeadID)
		}
	}

	return s.leadRepo.Update(ctx, lead)
}

func (s *LeadService) syncFromExternal(local *models.Lead, external *client.LeadResponse) {
	local.State = external.State
	local.Latitude = external.Latitude
	local.Longitude = external.Longitude
	local.Address = external.Address
	local.KwhUsage = external.KwhUsage
	local.SystemSize = external.SystemSize
	local.PanelCount = external.PanelCount
	local.PanelID = external.PanelID
	local.InverterID = external.InverterID
	local.UtilityID = external.UtilityID
	local.RoofMaterial = external.RoofMaterial
	local.AnnualProduction = external.AnnualProduction
	local.InstallationDate = external.InstallationDate
	local.MarkSynced(external.ID)
}

func (s *LeadService) syncExternalLeadToLocal(ctx context.Context, external *client.LeadResponse) error {
	leads, _ := s.leadRepo.ListByCompany(ctx, external.CompanyID, 1000, 0)
	
	var existingLead *models.Lead
	for _, lead := range leads {
		if lead.ExternalLeadID != nil && *lead.ExternalLeadID == external.ID {
			existingLead = lead
			break
		}
	}

	if existingLead != nil {
		s.syncFromExternal(existingLead, external)
		return s.leadRepo.Update(ctx, existingLead)
	}

	newLead := &models.Lead{
		CompanyID:        external.CompanyID,
		CreatorID:        external.CreatorID,
		State:            external.State,
		Latitude:         external.Latitude,
		Longitude:        external.Longitude,
		Address:          external.Address,
		Source:           external.Source,
		PromoCode:        external.PromoCode,
		Is2D:             external.Is2D,
		KwhUsage:         external.KwhUsage,
		SystemSize:       external.SystemSize,
		PanelCount:       external.PanelCount,
		PanelID:          external.PanelID,
		InverterID:       external.InverterID,
		UtilityID:        external.UtilityID,
		RoofMaterial:     external.RoofMaterial,
		AnnualProduction: external.AnnualProduction,
		InstallationDate: external.InstallationDate,
	}
	newLead.MarkSynced(external.ID)

	return s.leadRepo.Create(ctx, newLead)
}
