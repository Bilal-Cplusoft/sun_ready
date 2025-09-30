package models

import "time"

type Company struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	Name        string    `json:"name" gorm:"column:name"`
	DisplayName string    `json:"display_name" gorm:"column:display_name"`
	Description string    `json:"description" gorm:"column:description"`
	Code        string    `json:"code" gorm:"column:code"`
	Slug        string    `json:"slug" gorm:"column:slug;uniqueIndex"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true"`
	LogoPath    *string   `json:"logo_path" gorm:"column:logo_path"`
	AdminID     *int      `json:"admin_id" gorm:"column:admin_id"`
}

func (Company) TableName() string {
	return "companies"
}
