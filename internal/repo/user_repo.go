package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
	"errors"
)

var ErrUnauthorizedCompanyAccess = errors.New("insufficient permission to override company")

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


func (r *UserRepo) FindByIDs(ctx context.Context, ids []int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}


func (r *UserRepo) GetDescendantIDs(ctx context.Context, userID int) ([]int, error) {
	var descendants []int


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


func (r *UserRepo) UpdateCompanyID(ctx context.Context, userID, companyID int) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("company_id", companyID).Error
}

func (r *UserRepo) GetEffectiveCompanyID(ctx context.Context, user *models.User, companyID int) (int, error) {
	baseCompanyID := user.CompanyID
	if companyID == 0 {
		return baseCompanyID, nil
	}
	if user.Type == int16(models.UserTypeAdmin) || user.IsManager {
		return companyID, nil
	}

	return baseCompanyID, ErrUnauthorizedCompanyAccess
}
