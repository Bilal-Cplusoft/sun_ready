package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type LeadRepo struct {
	db *gorm.DB
}

func NewLeadRepo(db *gorm.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Create(ctx context.Context, lead *models.Lead) error {
	return r.db.WithContext(ctx).Create(lead).Error
}

func (r *LeadRepo) GetByID(ctx context.Context, id int) (*models.Lead, error) {
	var lead models.Lead
	err := r.db.WithContext(ctx).First(&lead, id).Error
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

func (r *LeadRepo) Update(ctx context.Context, lead *models.Lead) error {
	return r.db.WithContext(ctx).Save(lead).Error
}

func (r *LeadRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Lead{}, id).Error
}

func (r *LeadRepo) List(ctx context.Context, limit, offset int) ([]*models.Lead, error) {
	var leads []*models.Lead
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads).Error
	return leads, err
}

func (r *LeadRepo) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Lead, error) {
	var leads []*models.Lead
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads).Error
	return leads, err
}

func (r *LeadRepo) ListByCreator(ctx context.Context, creatorID int, limit, offset int) ([]*models.Lead, error) {
	var leads []*models.Lead
	err := r.db.WithContext(ctx).
		Where("creator_id = ?", creatorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads).Error
	return leads, err
}

func (r *LeadRepo) ListByState(ctx context.Context, state int, limit, offset int) ([]*models.Lead, error) {
	var leads []*models.Lead
	err := r.db.WithContext(ctx).
		Where("state = ?", state).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads).Error
	return leads, err
}

func (r *LeadRepo) UpdateState(ctx context.Context, id int, state int) error {
	return r.db.WithContext(ctx).Model(&models.Lead{}).
		Where("id = ?", id).
		Update("state", state).Error
}
