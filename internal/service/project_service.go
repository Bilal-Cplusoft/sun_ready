package service

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
)

type ProjectService struct {
	projectRepo *repo.ProjectRepo
}

func NewProjectService(projectRepo *repo.ProjectRepo) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(ctx context.Context, project *models.Project) error {
	return s.projectRepo.Create(ctx, project)
}

func (s *ProjectService) GetByID(ctx context.Context, id int) (*models.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *ProjectService) Update(ctx context.Context, project *models.Project) error {
	return s.projectRepo.Update(ctx, project)
}

func (s *ProjectService) Delete(ctx context.Context, id int) error {
	return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Project, error) {
	return s.projectRepo.ListByCompany(ctx, companyID, limit, offset)
}

func (s *ProjectService) ListByUser(ctx context.Context, userID int, limit, offset int) ([]*models.Project, error) {
	return s.projectRepo.ListByUser(ctx, userID, limit, offset)
}
