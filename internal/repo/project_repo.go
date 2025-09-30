package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepo) GetByID(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, id).Error
}

func (r *ProjectRepo) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID int, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error
	return projects, err
}
