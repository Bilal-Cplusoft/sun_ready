-- Migration: Add additional fields to companies table
-- This migration adds fields from the API project's company model

-- Add new columns to companies table
ALTER TABLE companies 
ADD COLUMN IF NOT EXISTS sales_commission_min DECIMAL(10,4),
ADD COLUMN IF NOT EXISTS sales_commission_max DECIMAL(10,4),
ADD COLUMN IF NOT EXISTS sales_commission_default DECIMAL(10,4),
ADD COLUMN IF NOT EXISTS baseline DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS baseline_adder DECIMAL(10,2),
ADD COLUMN IF NOT EXISTS baseline_adder_pct_sales_comms INTEGER,
ADD COLUMN IF NOT EXISTS contract_tag VARCHAR(100),
ADD COLUMN IF NOT EXISTS referred_by_user_id INTEGER,
ADD COLUMN IF NOT EXISTS credits INTEGER,
ADD COLUMN IF NOT EXISTS custom_commissions BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS pricing_mode INTEGER NOT NULL DEFAULT 0;

-- Add index for referred_by_user_id if needed
CREATE INDEX IF NOT EXISTS idx_companies_referred_by_user_id ON companies(referred_by_user_id);
