package models

import "time"

type Project struct {
	ID          int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	CompanyID   int       `json:"company_id" gorm:"column:company_id"`
	UserID      int       `json:"user_id" gorm:"column:user_id"`
	Name        string    `json:"name" gorm:"column:name"`
	Description *string   `json:"description" gorm:"column:description"`
	Status      string    `json:"status" gorm:"column:status;default:'draft'"`
	Address     *string   `json:"address" gorm:"column:address"`
}

func (Project) TableName() string {
	return "projects"
}

const (
	ProjectStatusDraft      = "draft"
	ProjectStatusInProgress = "in_progress"
	ProjectStatusCompleted  = "completed"
	ProjectStatusCancelled  = "cancelled"
)
