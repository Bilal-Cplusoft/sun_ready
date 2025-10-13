package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	FirstName   *string   `json:"first_name" gorm:"column:firstname"`
	LastName    *string   `json:"last_name" gorm:"column:lastname"`
	Email       string    `json:"email" gorm:"column:email;uniqueIndex"`
	Password    *string   `json:"-" gorm:"column:password"`
	Type        int16     `json:"type" gorm:"column:type"`
	PhoneNumber *string   `json:"phone_number" gorm:"column:phone_number"`
	Address     *string   `json:"address" gorm:"column:address"`
	CompanyID   int       `json:"company_id" gorm:"column:company_id"`
	CreatorID   *int      `json:"creator_id" gorm:"column:creator_id"`
	PicturePath *string   `json:"picture_path" gorm:"column:picture_path"`
	Disabled    bool      `json:"disabled" gorm:"column:disabled;default:false"`
	IsManager   bool      `json:"is_manager" gorm:"column:is_manager;default:false"`
}

func (User) TableName() string {
	return "users"
}

type UserType int16

const (
	UserTypeUnknown UserType = 0
	UserTypeAdmin   UserType = 1
	UserTypeSales   UserType = 2
	UserTypeClient  UserType = 3
)
