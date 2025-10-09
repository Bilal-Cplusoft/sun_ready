package models

import "errors"

var (
// Company errors
ErrInvalidCompanyName = errors.New("company name must be between 1 and 250 characters")
ErrInvalidCompanySlug = errors.New("company slug must be between 1 and 250 characters")
ErrCompanyNotFound    = errors.New("company not found")

// Deal errors
ErrInvalidDealTargetEPC        = errors.New("target EPC must be between 0 and 10000")
ErrInvalidDealHardwareCost     = errors.New("hardware cost must be between 0 and 10000000")
ErrInvalidDealInstallationCost = errors.New("installation cost must be between 0 and 10000000")
ErrInvalidDealSalesCommission  = errors.New("sales commission cost must be between 0 and 10000000")
ErrInvalidDealProfit           = errors.New("profit must be between 0 and 10000000")
ErrDealNotFound                = errors.New("deal not found")

// Lead errors
ErrInvalidLeadLatitude  = errors.New("latitude must be between -90 and 90")
ErrInvalidLeadLongitude = errors.New("longitude must be between -180 and 180")
ErrLeadNotFound         = errors.New("lead not found")

// Proposal errors
ErrInvalidProposalCode = errors.New("proposal code is required")
ErrInvalidProposalCost = errors.New("system cost must be greater than or equal to 0")
ErrProposalNotFound    = errors.New("proposal not found")

// Model3D errors
ErrInvalidModel3DLeadID      = errors.New("3D model must be associated with a valid lead")
ErrInvalidModel3DProjectID   = errors.New("3D model must have a valid LightFusion project ID")
ErrInvalidModel3DStatus      = errors.New("3D model status must be one of: pending, processing, completed, failed, expired")
ErrInvalidModel3DQuality     = errors.New("3D model quality must be one of: low, medium, high, ultra")
ErrInvalidModel3DRetryCount  = errors.New("3D model retry count must be between 0 and 10")
ErrInvalidModel3DSystemSize  = errors.New("3D model system size must be greater than or equal to 0")
ErrInvalidModel3DProduction  = errors.New("3D model annual production must be greater than or equal to 0")
ErrInvalidJSONData           = errors.New("invalid JSON data format")
ErrModel3DNotFound           = errors.New("3D model not found")
)
