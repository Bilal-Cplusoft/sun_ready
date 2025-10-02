package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type ProposalRepo struct {
	db *gorm.DB
}

func NewProposalRepo(db *gorm.DB) *ProposalRepo {
	return &ProposalRepo{db: db}
}

func (r *ProposalRepo) Create(ctx context.Context, proposal *models.Proposal) error {
	return r.db.WithContext(ctx).Create(proposal).Error
}

func (r *ProposalRepo) GetByID(ctx context.Context, id int) (*models.Proposal, error) {
	var proposal models.Proposal
	err := r.db.WithContext(ctx).First(&proposal, id).Error
	if err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (r *ProposalRepo) GetByCode(ctx context.Context, code string) (*models.Proposal, error) {
	var proposal models.Proposal
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&proposal).Error
	if err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (r *ProposalRepo) Update(ctx context.Context, proposal *models.Proposal) error {
	return r.db.WithContext(ctx).Save(proposal).Error
}

func (r *ProposalRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Proposal{}, id).Error
}

func (r *ProposalRepo) List(ctx context.Context, limit, offset int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&proposals).Error
	return proposals, err
}

func (r *ProposalRepo) ListByProject(ctx context.Context, projectID int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&proposals).Error
	return proposals, err
}

func (r *ProposalRepo) ListByCompany(ctx context.Context, companyID int, limit, offset int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&proposals).Error
	return proposals, err
}

func (r *ProposalRepo) ListBySales(ctx context.Context, salesID int, limit, offset int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Where("sales_id = ?", salesID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&proposals).Error
	return proposals, err
}

func (r *ProposalRepo) ListByHomeowner(ctx context.Context, homeownerID int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Where("homeowner_id = ?", homeownerID).
		Order("created_at DESC").
		Find(&proposals).Error
	return proposals, err
}

func (r *ProposalRepo) ListByStatus(ctx context.Context, status models.ProposalStatus, limit, offset int) ([]*models.Proposal, error) {
	var proposals []*models.Proposal
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&proposals).Error
	return proposals, err
}
