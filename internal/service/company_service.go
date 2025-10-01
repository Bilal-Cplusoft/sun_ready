package service

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
)

type CompanyService struct {
	companyRepo *repo.CompanyRepo
}

func NewCompanyService(companyRepo *repo.CompanyRepo) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) Create(ctx context.Context, company *models.Company) error {
	return s.companyRepo.Create(ctx, company)
}

func (s *CompanyService) GetByID(ctx context.Context, id int) (*models.Company, error) {
	return s.companyRepo.GetByID(ctx, id)
}

func (s *CompanyService) GetBySlug(ctx context.Context, slug string) (*models.Company, error) {
	return s.companyRepo.GetBySlug(ctx, slug)
}

func (s *CompanyService) Update(ctx context.Context, company *models.Company) error {
	return s.companyRepo.Update(ctx, company)
}

func (s *CompanyService) List(ctx context.Context, limit, offset int) ([]*models.Company, error) {
	return s.companyRepo.List(ctx, limit, offset)
}

func (s *CompanyService) FindAll(ctx context.Context) ([]*models.Company, error) {
	return s.companyRepo.FindAll(ctx)
}

func (s *CompanyService) Delete(ctx context.Context, id int) error {
	return s.companyRepo.Delete(ctx, id)
}
