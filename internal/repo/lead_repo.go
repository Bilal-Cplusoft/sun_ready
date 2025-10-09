package repo

import (
	"context"
	"fmt"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

// LeadRepo handles lead database operations
type LeadRepo struct {
	db *gorm.DB
}

// NewLeadRepo creates a new lead repository
func NewLeadRepo(db *gorm.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

// Create creates a new lead in the database
func (r *LeadRepo) Create(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Create(lead)
	if result.Error != nil {
		return fmt.Errorf("failed to create lead: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a lead by ID
func (r *LeadRepo) GetByID(ctx context.Context, id int) (*models.Lead, error) {
	var lead models.Lead
	result := r.db.WithContext(ctx).First(&lead, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrLeadNotFound
		}
		return nil, fmt.Errorf("failed to get lead: %w", result.Error)
	}

	return &lead, nil
}

// GetByExternalID retrieves a lead by external (LightFusion) ID
func (r *LeadRepo) GetByExternalID(ctx context.Context, externalID int) (*models.Lead, error) {
	var lead models.Lead
	result := r.db.WithContext(ctx).Where("external_lead_id = ?", externalID).First(&lead)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrLeadNotFound
		}
		return nil, fmt.Errorf("failed to get lead by external ID: %w", result.Error)
	}

	return &lead, nil
}

// GetByLightFusionProjectID retrieves a lead by LightFusion 3D project ID
func (r *LeadRepo) GetByLightFusionProjectID(ctx context.Context, projectID int) (*models.Lead, error) {
	var lead models.Lead
	result := r.db.WithContext(ctx).Where("lightfusion_3d_project_id = ?", projectID).First(&lead)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrLeadNotFound
		}
		return nil, fmt.Errorf("failed to get lead by LightFusion project ID: %w", result.Error)
	}

	return &lead, nil
}

// Update updates an existing lead
func (r *LeadRepo) Update(ctx context.Context, lead *models.Lead) error {
	if err := lead.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result := r.db.WithContext(ctx).Save(lead)
	if result.Error != nil {
		return fmt.Errorf("failed to update lead: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}

// Delete deletes a lead by ID
func (r *LeadRepo) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Lead{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete lead: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}

// List retrieves all leads with optional filters
func (r *LeadRepo) List(ctx context.Context, companyID *int, creatorID *int, limit, offset int) ([]*models.Lead, int64, error) {
	var leads []*models.Lead
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Lead{})

	// Apply filters
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}
	if creatorID != nil {
		query = query.Where("creator_id = ?", *creatorID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count leads: %w", err)
	}

	// Get paginated results
	result := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", result.Error)
	}

	return leads, total, nil
}

// ListWith3DModels retrieves leads that have 3D models
func (r *LeadRepo) ListWith3DModels(ctx context.Context, companyID *int, limit, offset int) ([]*models.Lead, int64, error) {
	var leads []*models.Lead
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Lead{}).
		Where("lightfusion_3d_project_id IS NOT NULL")

	// Apply company filter if provided
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count leads with 3D models: %w", err)
	}

	// Get paginated results
	result := query.
		Order("model_3d_created_at DESC NULLS LAST, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&leads)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list leads with 3D models: %w", result.Error)
	}

	return leads, total, nil
}

// ListPendingSync retrieves leads that need to be synced with external API
func (r *LeadRepo) ListPendingSync(ctx context.Context, limit int) ([]*models.Lead, error) {
	var leads []*models.Lead

	result := r.db.WithContext(ctx).
		Where("sync_status IN ?", []string{"pending", "failed"}).
		Order("created_at ASC").
		Limit(limit).
		Find(&leads)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list pending sync leads: %w", result.Error)
	}

	return leads, nil
}

// UpdateSyncStatus updates the sync status of a lead
func (r *LeadRepo) UpdateSyncStatus(ctx context.Context, leadID int, status string, externalID *int) error {
	updates := map[string]interface{}{
		"sync_status": status,
	}
	
	if externalID != nil {
		updates["external_lead_id"] = *externalID
		updates["last_synced_at"] = gorm.Expr("NOW()")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Lead{}).
		Where("id = ?", leadID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update sync status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}

// Update3DModelStatus updates the 3D model status of a lead
func (r *LeadRepo) Update3DModelStatus(ctx context.Context, leadID int, status string) error {
	updates := map[string]interface{}{
		"model_3d_status": status,
	}
	
	if status == "completed" {
		updates["model_3d_completed_at"] = gorm.Expr("NOW()")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Lead{}).
		Where("id = ?", leadID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update 3D model status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrLeadNotFound
	}

	return nil
}

// GetLeadsWith3DModelsByStatus retrieves leads with specific 3D model status
func (r *LeadRepo) GetLeadsWith3DModelsByStatus(ctx context.Context, status string, limit int) ([]*models.Lead, error) {
	var leads []*models.Lead

	result := r.db.WithContext(ctx).
		Where("model_3d_status = ?", status).
		Order("model_3d_created_at DESC").
		Limit(limit).
		Find(&leads)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get leads by 3D model status: %w", result.Error)
	}

	return leads, nil
}

// BatchUpdate3DModelStatus updates multiple leads' 3D model status
func (r *LeadRepo) BatchUpdate3DModelStatus(ctx context.Context, leadIDs []int, status string) error {
	if len(leadIDs) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"model_3d_status": status,
		"updated_at": gorm.Expr("NOW()"),
	}
	
	if status == "completed" {
		updates["model_3d_completed_at"] = gorm.Expr("NOW()")
	}

	result := r.db.WithContext(ctx).
		Model(&models.Lead{}).
		Where("id IN ?", leadIDs).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to batch update 3D model status: %w", result.Error)
	}

	return nil
}