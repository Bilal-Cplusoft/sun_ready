-- Migration: Create deals table
-- This migration creates the deals table for managing solar installation deals

CREATE TABLE IF NOT EXISTS deals (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    uuid VARCHAR(36) NOT NULL UNIQUE,
    signed_at TIMESTAMPTZ,
    
    -- Foreign Keys
    lead_id INTEGER,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    system_id INTEGER,
    hardware_id INTEGER,
    sales_id INTEGER NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    homeowner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    document_id INTEGER,
    financing_option_id INTEGER,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE NO ACTION,
    
    -- System Details
    system_size DECIMAL(10,2) NOT NULL DEFAULT 0,
    panel_count INTEGER NOT NULL DEFAULT 0,
    panel_id INTEGER,
    inverter_id INTEGER,
    
    -- Financial Details
    financing_provider VARCHAR(255),
    target_epc DECIMAL(10,4) NOT NULL,
    total_cost DECIMAL(12,2) NOT NULL,
    hardware_cost DECIMAL(12,2) NOT NULL,
    installation_cost DECIMAL(12,2) NOT NULL,
    sales_commission_cost DECIMAL(12,2) NOT NULL,
    profit DECIMAL(12,2) NOT NULL,
    
    -- Status and Metadata
    archive BOOLEAN NOT NULL DEFAULT false,
    approved_at TIMESTAMPTZ,
    installed_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    
    -- Additional Info
    address TEXT,
    consumption_kwh INTEGER DEFAULT 0,
    production_kwh INTEGER DEFAULT 0,
    utility_id INTEGER
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_deals_uuid ON deals(uuid);
CREATE INDEX IF NOT EXISTS idx_deals_lead_id ON deals(lead_id);
CREATE INDEX IF NOT EXISTS idx_deals_project_id ON deals(project_id);
CREATE INDEX IF NOT EXISTS idx_deals_system_id ON deals(system_id);
CREATE INDEX IF NOT EXISTS idx_deals_hardware_id ON deals(hardware_id);
CREATE INDEX IF NOT EXISTS idx_deals_sales_id ON deals(sales_id);
CREATE INDEX IF NOT EXISTS idx_deals_homeowner_id ON deals(homeowner_id);
CREATE INDEX IF NOT EXISTS idx_deals_company_id ON deals(company_id);
CREATE INDEX IF NOT EXISTS idx_deals_signed_at ON deals(signed_at);
CREATE INDEX IF NOT EXISTS idx_deals_status ON deals(status);
CREATE INDEX IF NOT EXISTS idx_deals_archive ON deals(archive);
CREATE INDEX IF NOT EXISTS idx_deals_panel_id ON deals(panel_id);
CREATE INDEX IF NOT EXISTS idx_deals_inverter_id ON deals(inverter_id);
