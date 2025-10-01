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
)
