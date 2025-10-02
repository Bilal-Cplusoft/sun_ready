package models

import (
	"time"
)

// ProposalStatus represents the status of a proposal
type ProposalStatus string

const (
	ProposalStatusDraft     ProposalStatus = "draft"
	ProposalStatusSent      ProposalStatus = "sent"
	ProposalStatusViewed    ProposalStatus = "viewed"
	ProposalStatusAccepted  ProposalStatus = "accepted"
	ProposalStatusRejected  ProposalStatus = "rejected"
	ProposalStatusExpired   ProposalStatus = "expired"
)

type Proposal struct {
	ID          int            `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
	Code        string         `json:"code" gorm:"column:code;uniqueIndex;not null" example:"PROP-ABC123"`
	Status      ProposalStatus `json:"status" gorm:"column:status;not null;default:'draft'" example:"draft"`
	
	// Relationships
	ProjectID   int            `json:"project_id" gorm:"column:project_id;not null" example:"1"`
	LeadID      *int           `json:"lead_id" gorm:"column:lead_id" example:"1"`
	CompanyID   int            `json:"company_id" gorm:"column:company_id;not null" example:"1"`
	SalesID     int            `json:"sales_id" gorm:"column:sales_id;not null" example:"1"`
	HomeownerID int            `json:"homeowner_id" gorm:"column:homeowner_id;not null" example:"2"`
	
	// System Details
	SystemSize       float64 `json:"system_size" gorm:"column:system_size" example:"10.5"`
	PanelCount       int     `json:"panel_count" gorm:"column:panel_count" example:"30"`
	PanelID          *int    `json:"panel_id" gorm:"column:panel_id" example:"1"`
	InverterID       *int    `json:"inverter_id" gorm:"column:inverter_id" example:"1"`
	BatteryCount     int     `json:"battery_count" gorm:"column:battery_count;default:0" example:"0"`
	
	// Production & Consumption
	AnnualProduction  float64 `json:"annual_production" gorm:"column:annual_production" example:"13000"`
	AnnualConsumption float64 `json:"annual_consumption" gorm:"column:annual_consumption" example:"12000"`
	
	// Financial Details
	SystemCost           float64 `json:"system_cost" gorm:"column:system_cost" example:"25000.00"`
	Incentives           float64 `json:"incentives" gorm:"column:incentives" example:"5000.00"`
	NetCost              float64 `json:"net_cost" gorm:"column:net_cost" example:"20000.00"`
	MonthlyPayment       float64 `json:"monthly_payment" gorm:"column:monthly_payment" example:"150.00"`
	FinancingOptionID    *int    `json:"financing_option_id" gorm:"column:financing_option_id" example:"1"`
	FinancingProvider    string  `json:"financing_provider" gorm:"column:financing_provider" example:"SunPower Financial"`
	
	// Utility Details
	UtilityID            *int    `json:"utility_id" gorm:"column:utility_id" example:"1"`
	CurrentUtilityBill   float64 `json:"current_utility_bill" gorm:"column:current_utility_bill" example:"200.00"`
	EstimatedUtilityBill float64 `json:"estimated_utility_bill" gorm:"column:estimated_utility_bill" example:"50.00"`
	
	// Document Details
	DocumentID      *int       `json:"document_id" gorm:"column:document_id" example:"1"`
	DocumentURL     string     `json:"document_url" gorm:"column:document_url" example:"https://docs.example.com/proposal.pdf"`
	ContractURL     string     `json:"contract_url" gorm:"column:contract_url" example:"https://docs.example.com/contract.pdf"`
	
	// Timestamps
	SentAt     *time.Time `json:"sent_at" gorm:"column:sent_at" example:"2025-10-01T10:00:00Z"`
	ViewedAt   *time.Time `json:"viewed_at" gorm:"column:viewed_at" example:"2025-10-01T11:00:00Z"`
	AcceptedAt *time.Time `json:"accepted_at" gorm:"column:accepted_at" example:"2025-10-01T12:00:00Z"`
	ExpiresAt  *time.Time `json:"expires_at" gorm:"column:expires_at" example:"2025-10-30T23:59:59Z"`
	
	// Additional Info
	Notes      string `json:"notes" gorm:"column:notes;type:text" example:"Custom proposal notes"`
	Address    string `json:"address" gorm:"column:address" example:"123 Solar St, San Francisco, CA 94102"`
}

func (Proposal) TableName() string {
	return "proposals"
}

// Validate validates proposal data
func (p *Proposal) Validate() error {
	if p.Code == "" {
		return ErrInvalidProposalCode
	}
	if p.SystemCost < 0 {
		return ErrInvalidProposalCost
	}
	return nil
}

// IsExpired checks if the proposal has expired
func (p *Proposal) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

// MarkSent marks the proposal as sent
func (p *Proposal) MarkSent() {
	p.Status = ProposalStatusSent
	now := time.Now()
	p.SentAt = &now
	
	// Set expiration to 30 days from now if not set
	if p.ExpiresAt == nil {
		expiresAt := now.AddDate(0, 0, 30)
		p.ExpiresAt = &expiresAt
	}
}

// MarkViewed marks the proposal as viewed
func (p *Proposal) MarkViewed() {
	if p.Status == ProposalStatusDraft {
		return
	}
	p.Status = ProposalStatusViewed
	if p.ViewedAt == nil {
		now := time.Now()
		p.ViewedAt = &now
	}
}

// MarkAccepted marks the proposal as accepted
func (p *Proposal) MarkAccepted() {
	p.Status = ProposalStatusAccepted
	now := time.Now()
	p.AcceptedAt = &now
}

// MarkRejected marks the proposal as rejected
func (p *Proposal) MarkRejected() {
	p.Status = ProposalStatusRejected
}
