package models

import (
	"strings"
	"time"
)

type Company struct {
	ID                         int       `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt                  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                  time.Time `json:"updated_at" gorm:"column:updated_at"`
	Name                       string    `json:"name" gorm:"column:name;not null" example:"Acme Corp"`
	DisplayName                string    `json:"display_name" gorm:"column:display_name" example:"Acme Corporation"`
	Description                string    `json:"description" gorm:"column:description" example:"Leading solar company"`
	Code                       string    `json:"code" gorm:"column:code" example:"ACME"`
	Slug                       string    `json:"slug" gorm:"column:slug;uniqueIndex;not null" example:"acme-corp"`
	IsActive                   bool      `json:"is_active" gorm:"column:is_active;default:true" example:"true"`
	LogoPath                   *string   `json:"logo_path" gorm:"column:logo_path" example:"https://example.com/logo.png"`
	AdminID                    *int      `json:"admin_id" gorm:"column:admin_id" example:"1"`
	SalesCommissionMin         *float64  `json:"sales_commission_min" gorm:"column:sales_commission_min" example:"0.05"`
	SalesCommissionMax         *float64  `json:"sales_commission_max" gorm:"column:sales_commission_max" example:"0.15"`
	SalesCommissionDefault     *float64  `json:"sales_commission_default" gorm:"column:sales_commission_default" example:"0.10"`
	Baseline                   *float64  `json:"baseline" gorm:"column:baseline" example:"1000.00"`
	BaselineAdder              *float64  `json:"baseline_adder" gorm:"column:baseline_adder" example:"100.00"`
	BaselineAdderPctSalesComms *int      `json:"baseline_adder_pct_sales_comms" gorm:"column:baseline_adder_pct_sales_comms" example:"10"`
	ContractTag                *string   `json:"contract_tag" gorm:"column:contract_tag" example:"STANDARD"`
	ReferredByUserID           *int      `json:"referred_by_user_id" gorm:"column:referred_by_user_id" example:"1"`
	Credits                    *int      `json:"credits" gorm:"column:credits" example:"1000"`
	CustomCommissions          bool      `json:"custom_commissions" gorm:"column:custom_commissions;default:false" example:"false"`
	PricingMode                int       `json:"pricing_mode" gorm:"column:pricing_mode;default:0" example:"0"`
}

func (Company) TableName() string {
	return "companies"
}

// Sanitize will sanitize company name
func (c *Company) Sanitize() {
	c.Name = strings.TrimSpace(c.Name)
	c.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(c.Slug), " ", "-"))
}

// Validate will validate company
func (c *Company) Validate() error {
	if len(c.Name) == 0 || len(c.Name) > 250 {
		return ErrInvalidCompanyName
	}
	if len(c.Slug) == 0 || len(c.Slug) > 250 {
		return ErrInvalidCompanySlug
	}
	return nil
}
