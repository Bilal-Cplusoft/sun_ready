package repo

import (
	"context"

	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"gorm.io/gorm"
)

type QuoteRepo struct {
	db *gorm.DB
}

func NewQuoteRepo(db *gorm.DB) *QuoteRepo {
	return &QuoteRepo{db: db}
}

func (r *QuoteRepo) Create(ctx context.Context, deal *models.Deal) error {
	return r.db.WithContext(ctx).Create(deal).Error
}
