package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepo) List(ctx context.Context, companyID int, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}

// FindByIDs finds users by their IDs
func (r *UserRepo) FindByIDs(ctx context.Context, ids []int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}

// GetDescendantIDs gets all descendant user IDs for a given user
func (r *UserRepo) GetDescendantIDs(ctx context.Context, userID int) ([]int, error) {
	var descendants []int
	
	// Recursive CTE to get all descendants
	query := `
		WITH RECURSIVE user_tree AS (
			SELECT id, creator_id FROM users WHERE creator_id = ?
			UNION ALL
			SELECT u.id, u.creator_id FROM users u
			INNER JOIN user_tree ut ON u.creator_id = ut.id
		)
		SELECT id FROM user_tree
	`
	
	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&descendants).Error
	return descendants, err
}

// UpdateCompanyID updates the company_id for a user
func (r *UserRepo) UpdateCompanyID(ctx context.Context, userID, companyID int) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("company_id", companyID).Error
}
