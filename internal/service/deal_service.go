package service

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
)

type DealService struct {
	dealRepo *repo.DealRepo
}

func NewDealService(dealRepo *repo.DealRepo) *DealService {
	return &DealService{dealRepo: dealRepo}
}

func (s *DealService) Create(ctx context.Context, deal *models.Deal) error {
	if err := deal.Validate(); err != nil {
		return err
	}
	return s.dealRepo.Create(ctx, deal)
}

func (s *DealService) GetByID(ctx context.Context, id int) (*models.Deal, error) {
	return s.dealRepo.GetByID(ctx, id)
}

func (s *DealService) GetByUUID(ctx context.Context, uuid string) (*models.Deal, error) {
	return s.dealRepo.GetByUUID(ctx, uuid)
}

func (s *DealService) Update(ctx context.Context, deal *models.Deal) error {
	if err := deal.Validate(); err != nil {
		return err
	}
	return s.dealRepo.Update(ctx, deal)
}

func (s *DealService) Delete(ctx context.Context, id int) error {
	return s.dealRepo.Delete(ctx, id)
}

func (s *DealService) List(ctx context.Context, limit, offset int) ([]*models.Deal, error) {
	return s.dealRepo.List(ctx, limit, offset)
}

func (s *DealService) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Deal, error) {
	return s.dealRepo.ListByCompany(ctx, companyID, limit, offset)
}

func (s *DealService) ListBySales(ctx context.Context, salesID int, limit, offset int) ([]*models.Deal, error) {
	return s.dealRepo.ListBySales(ctx, salesID, limit, offset)
}

func (s *DealService) ListByHomeowner(ctx context.Context, homeownerID int, limit, offset int) ([]*models.Deal, error) {
	return s.dealRepo.ListByHomeowner(ctx, homeownerID, limit, offset)
}

func (s *DealService) ListByProject(ctx context.Context, projectID int) ([]*models.Deal, error) {
	return s.dealRepo.ListByProject(ctx, projectID)
}

func (s *DealService) ListSigned(ctx context.Context, companyID int, limit, offset int) ([]*models.Deal, error) {
	return s.dealRepo.ListSigned(ctx, companyID, limit, offset)
}

func (s *DealService) Archive(ctx context.Context, id int) error {
	return s.dealRepo.Archive(ctx, id)
}

func (s *DealService) Unarchive(ctx context.Context, id int) error {
	return s.dealRepo.Unarchive(ctx, id)
}
