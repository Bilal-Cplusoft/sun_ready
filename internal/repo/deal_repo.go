package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type DealRepo struct {
	db *gorm.DB
}

func NewDealRepo(db *gorm.DB) *DealRepo {
	return &DealRepo{db: db}
}

func (r *DealRepo) Create(ctx context.Context, deal *models.Deal) error {
	return r.db.WithContext(ctx).Create(deal).Error
}

func (r *DealRepo) GetByID(ctx context.Context, id int) (*models.Deal, error) {
	var deal models.Deal
	err := r.db.WithContext(ctx).First(&deal, id).Error
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealRepo) GetByUUID(ctx context.Context, uuid string) (*models.Deal, error) {
	var deal models.Deal
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&deal).Error
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealRepo) Update(ctx context.Context, deal *models.Deal) error {
	return r.db.WithContext(ctx).Save(deal).Error
}

func (r *DealRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Deal{}, id).Error
}

func (r *DealRepo) List(ctx context.Context, limit, offset int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) ListBySales(ctx context.Context, salesID int, limit, offset int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Where("sales_id = ?", salesID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) ListByHomeowner(ctx context.Context, homeownerID int, limit, offset int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Where("homeowner_id = ?", homeownerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) ListByProject(ctx context.Context, projectID int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) ListSigned(ctx context.Context, companyID int, limit, offset int) ([]*models.Deal, error) {
	var deals []*models.Deal
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND signed_at IS NOT NULL", companyID).
		Order("signed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&deals).Error
	return deals, err
}

func (r *DealRepo) Archive(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&models.Deal{}).
		Where("id = ?", id).
		Update("archive", true).Error
}

func (r *DealRepo) Unarchive(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&models.Deal{}).
		Where("id = ?", id).
		Update("archive", false).Error
}
