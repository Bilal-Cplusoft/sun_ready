package models

import (
"time"
)

type Deal struct {
ID                   int        `json:"id" gorm:"primaryKey;column:id"`
CreatedAt            time.Time  `json:"created_at" gorm:"column:created_at"`
UpdatedAt            time.Time  `json:"updated_at" gorm:"column:updated_at"`
UUID                 string     `json:"uuid" gorm:"column:uuid;uniqueIndex;not null" example:"550e8400-e29b-41d4-a716-446655440000"`
SignedAt             *time.Time `json:"signed_at" gorm:"column:signed_at" example:"2025-10-01T10:00:00Z"`

// Foreign Keys
LeadID               *int       `json:"lead_id" gorm:"column:lead_id" example:"1"`
ProjectID            int        `json:"project_id" gorm:"column:project_id;not null" example:"1"`
SystemID             *int       `json:"system_id" gorm:"column:system_id" example:"1"`
HardwareID           *int       `json:"hardware_id" gorm:"column:hardware_id" example:"1"`
SalesID              int        `json:"sales_id" gorm:"column:sales_id;not null" example:"1"`
HomeownerID          int        `json:"homeowner_id" gorm:"column:homeowner_id;not null" example:"2"`
DocumentID           *int       `json:"document_id" gorm:"column:document_id" example:"1"`
FinancingOptionID    *int       `json:"financing_option_id" gorm:"column:financing_option_id" example:"1"`
CompanyID            int        `json:"company_id" gorm:"column:company_id;not null" example:"1"`

// System Details
SystemSize           float64    `json:"system_size" gorm:"column:system_size" example:"10.5"`
PanelCount           int        `json:"panel_count" gorm:"column:panel_count" example:"30"`
PanelID              *int       `json:"panel_id" gorm:"column:panel_id" example:"1"`
InverterID           *int       `json:"inverter_id" gorm:"column:inverter_id" example:"1"`

// Financial Details
FinancingProvider    string     `json:"financing_provider" gorm:"column:financing_provider" example:"SunPower Financial"`
TargetEPC            float64    `json:"target_epc" gorm:"column:target_epc;not null" example:"2.50"`
TotalCost            float64    `json:"total_cost" gorm:"column:total_cost;not null" example:"25000.00"`
HardwareCost         float64    `json:"hardware_cost" gorm:"column:hardware_cost;not null" example:"15000.00"`
InstallationCost     float64    `json:"installation_cost" gorm:"column:installation_cost;not null" example:"8000.00"`
SalesCommissionCost  float64    `json:"sales_commission_cost" gorm:"column:sales_commission_cost;not null" example:"2000.00"`
Profit               float64    `json:"profit" gorm:"column:profit;not null" example:"5000.00"`

// Status and Metadata
Archive              bool       `json:"archive" gorm:"column:archive;default:false" example:"false"`
ApprovedAt           *time.Time `json:"approved_at" gorm:"column:approved_at" example:"2025-10-02T10:00:00Z"`
InstalledAt          *time.Time `json:"installed_at" gorm:"column:installed_at" example:"2025-10-15T10:00:00Z"`
Status               string     `json:"status" gorm:"column:status;default:'pending'" example:"pending"`

// Additional Info
Address              string     `json:"address" gorm:"column:address" example:"123 Solar Street, CA 90210"`
ConsumptionKWH       int        `json:"consumption_kwh" gorm:"column:consumption_kwh" example:"12000"`
ProductionKWH        int        `json:"production_kwh" gorm:"column:production_kwh" example:"13000"`
UtilityID            *int       `json:"utility_id" gorm:"column:utility_id" example:"1"`
}

func (Deal) TableName() string {
return "deals"
}

// Validate validates deal data
func (d *Deal) Validate() error {
if d.TargetEPC < 0 || d.TargetEPC > 10000 {
		return ErrInvalidDealTargetEPC
}
if d.HardwareCost < 0 || d.HardwareCost > 10000000 {
		return ErrInvalidDealHardwareCost
}
if d.InstallationCost < 0 || d.InstallationCost > 10000000 {
		return ErrInvalidDealInstallationCost
}
if d.SalesCommissionCost < 0 || d.SalesCommissionCost > 10000000 {
		return ErrInvalidDealSalesCommission
}
if d.Profit < 0 || d.Profit > 10000000 {
		return ErrInvalidDealProfit
}
return nil
}
