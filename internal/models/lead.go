package models

import (
"time"
)

// LeadSource represents the source of the lead
type LeadSource int

const (
LeadSourceLegacy  LeadSource = 0
LeadSourceDrone   LeadSource = 1
LeadSourceEarth   LeadSource = 2
LeadSourceFlyover LeadSource = 3
LeadSourceNone    LeadSource = 4
)

// LeadState represents the state of the lead
type LeadState int

const (
LeadStateProgress     LeadState = 0
LeadStateDone         LeadState = 1
LeadStateErrored      LeadState = 2
LeadStateInitialized  LeadState = 3
)

// Lead represents a local reference to a lead (stored in both local DB and external API)
type Lead struct {
ID                  int        `json:"id" gorm:"primaryKey;column:id"`
CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at"`
UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at"`

// External reference
ExternalLeadID      *int       `json:"external_lead_id" gorm:"column:external_lead_id;uniqueIndex" example:"12345"`
SyncStatus          string     `json:"sync_status" gorm:"column:sync_status;default:'pending'" example:"synced"`
LastSyncedAt        *time.Time `json:"last_synced_at" gorm:"column:last_synced_at"`

// Core fields (synced with external API)
State               int        `json:"state" gorm:"column:state;not null;default:0" example:"0"`
CompanyID           int        `json:"company_id" gorm:"column:company_id;not null" example:"1"`
CreatorID           int        `json:"creator_id" gorm:"column:creator_id;not null" example:"1"`

// Location
Latitude            float64    `json:"latitude" gorm:"column:latitude;not null" example:"37.7749"`
Longitude           float64    `json:"longitude" gorm:"column:longitude;not null" example:"-122.4194"`
Address             string     `json:"address" gorm:"column:address" example:"123 Solar St, San Francisco, CA 94102"`

// Source and metadata
Source              int        `json:"source" gorm:"column:source;not null;default:0" example:"0"`
PromoCode           *string    `json:"promo_code" gorm:"column:promo_code" example:"SOLAR2025"`
Is2D                bool       `json:"is_2d" gorm:"column:is_2d;default:false" example:"false"`

// Energy details
KwhUsage            float64    `json:"kwh_usage" gorm:"column:kwh_usage" example:"12000"`
KwhPerKwManual      int        `json:"kwh_per_kw_manual" gorm:"column:kwh_per_kw_manual" example:"1200"`

// Financial
ElectricityCostPre  *int       `json:"electricity_cost_pre" gorm:"column:electricity_cost_pre" example:"150"`
ElectricityCostPost *int       `json:"electricity_cost_post" gorm:"column:electricity_cost_post" example:"50"`
AdditionalIncentive *int       `json:"additional_incentive" gorm:"column:additional_incentive" example:"1000"`

// System details
SystemSize          float64    `json:"system_size" gorm:"column:system_size" example:"10.5"`
PanelCount          int        `json:"panel_count" gorm:"column:panel_count" example:"30"`
PanelID             *int       `json:"panel_id" gorm:"column:panel_id" example:"1"`
InverterID          *int       `json:"inverter_id" gorm:"column:inverter_id" example:"1"`
InverterCount       int        `json:"inverter_count" gorm:"column:inverter_count;default:1" example:"1"`
BatteryCount        int        `json:"battery_count" gorm:"column:battery_count;default:0" example:"0"`

// Utility
UtilityID           *int       `json:"utility_id" gorm:"column:utility_id" example:"1"`
TariffID            *int       `json:"tariff_id" gorm:"column:tariff_id" example:"1"`

// Roof details
RoofMaterial        *int       `json:"roof_material" gorm:"column:roof_material" example:"1"`
SurfaceID           *int       `json:"surface_id" gorm:"column:surface_id" example:"1"`

// Production
AnnualProduction    float64    `json:"annual_production" gorm:"column:annual_production" example:"13000"`

// Workflow states
WelcomeCallState       *int `json:"welcome_call_state" gorm:"column:welcome_call_state" example:"0"`
FinancingState         *int `json:"financing_state" gorm:"column:financing_state" example:"0"`
UtilityBillState       *int `json:"utility_bill_state" gorm:"column:utility_bill_state" example:"0"`
DesignApprovedState    *int `json:"design_approved_state" gorm:"column:design_approved_state" example:"0"`
PermittingApprovedState *int `json:"permitting_approved_state" gorm:"column:permitting_approved_state" example:"0"`
SitePhotosState        *int `json:"site_photos_state" gorm:"column:site_photos_state" example:"0"`
InstallCrewState       *int `json:"install_crew_state" gorm:"column:install_crew_state" example:"0"`
InstallationState      *int `json:"installation_state" gorm:"column:installation_state" example:"0"`
FinalInspectionState   *int `json:"final_inspection_state" gorm:"column:final_inspection_state" example:"0"`
PtoState               *int `json:"pto_state" gorm:"column:pto_state" example:"0"`

// Dates
InstallationDate *string `json:"installation_date" gorm:"column:installation_date" example:"2025-10-15"`
DateNtp          *string `json:"date_ntp" gorm:"column:date_ntp" example:"2025-10-10"`
DateInstalled    *string `json:"date_installed" gorm:"column:date_installed" example:"2025-10-15"`

// LightFusion 3D Project Integration
LightFusion3DProjectID *int    `json:"lightfusion_3d_project_id" gorm:"column:lightfusion_3d_project_id" example:"123"`
LightFusion3DHouseID   *int    `json:"lightfusion_3d_house_id" gorm:"column:lightfusion_3d_house_id" example:"456"`
Model3DStatus          *string `json:"model_3d_status" gorm:"column:model_3d_status" example:"processing"`
Model3DCreatedAt       *time.Time `json:"model_3d_created_at" gorm:"column:model_3d_created_at"`
Model3DCompletedAt     *time.Time `json:"model_3d_completed_at" gorm:"column:model_3d_completed_at"`
}

func (Lead) TableName() string {
return "leads"
}

// Validate validates lead data
func (l *Lead) Validate() error {
if l.Latitude < -90 || l.Latitude > 90 {
		return ErrInvalidLeadLatitude
}
if l.Longitude < -180 || l.Longitude > 180 {
		return ErrInvalidLeadLongitude
}
return nil
}

func (l *Lead) MarkSynced(externalID int) {
l.ExternalLeadID = &externalID
l.SyncStatus = "synced"
now := time.Now()
l.LastSyncedAt = &now
}

func (l *Lead) MarkSyncFailed() {
l.SyncStatus = "failed"
}


func (l *Lead) Has3DModel() bool {
	return l.LightFusion3DProjectID != nil && l.LightFusion3DHouseID != nil
}

func (l *Lead) SetLightFusion3DProject(projectID, houseID int) {
	l.LightFusion3DProjectID = &projectID
	l.LightFusion3DHouseID = &houseID
	status := "processing"
	l.Model3DStatus = &status
	now := time.Now()
	l.Model3DCreatedAt = &now
}

func (l *Lead) Update3DModelStatus(status string) {
	l.Model3DStatus = &status
	if status == "completed" || status == "done" {
		now := time.Now()
		l.Model3DCompletedAt = &now
	}
}

func (l *Lead) Is3DModelReady() bool {
	return l.Model3DStatus != nil && (*l.Model3DStatus == "completed" || *l.Model3DStatus == "done")
}
