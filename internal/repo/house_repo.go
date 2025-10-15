package repo

import (
	"gorm.io/gorm"
	"context"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
)


type HouseRepo struct {
	db *gorm.DB
}


func NewHouseRepo(db *gorm.DB) *HouseRepo {
	return &HouseRepo{db: db}
}

func (r *HouseRepo) Create(ctx context.Context, house *models.House) error {
	if err := r.db.WithContext(ctx).Create(house).Error; err != nil {
		return err
	}
	return nil
}
