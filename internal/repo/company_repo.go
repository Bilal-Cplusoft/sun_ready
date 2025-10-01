package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type CompanyRepo struct {
	db *gorm.DB
}

func NewCompanyRepo(db *gorm.DB) *CompanyRepo {
	return &CompanyRepo{db: db}
}

func (r *CompanyRepo) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *CompanyRepo) GetByID(ctx context.Context, id int) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepo) GetBySlug(ctx context.Context, slug string) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepo) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *CompanyRepo) List(ctx context.Context, limit, offset int) ([]*models.Company, error) {
	var companies []*models.Company
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&companies).Error
	return companies, err
}

func (r *CompanyRepo) FindAll(ctx context.Context) ([]*models.Company, error) {
	var companies []*models.Company
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&companies).Error
	return companies, err
}

func (r *CompanyRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Company{}, id).Error
}
