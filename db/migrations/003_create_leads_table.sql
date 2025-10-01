-- Migration: Create leads table with external API sync support
-- This migration creates the leads table for managing solar leads with LightFUSION API integration

CREATE TABLE IF NOT EXISTS leads (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- External API sync fields
    external_lead_id INTEGER UNIQUE,
    sync_status VARCHAR(20) DEFAULT 'pending',
    last_synced_at TIMESTAMPTZ,
    
    -- Core fields
    state INTEGER NOT NULL DEFAULT 0,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE NO ACTION,
    creator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    
    -- Location
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    address TEXT,
    
    -- Source and metadata
    source INTEGER NOT NULL DEFAULT 0,
    promo_code VARCHAR(50),
    is_2d BOOLEAN DEFAULT false,
    
    -- Energy details
    kwh_usage DECIMAL(12,2) DEFAULT 0,
    kwh_per_kw_manual INTEGER DEFAULT 0,
    
    -- Financial
    electricity_cost_pre INTEGER,
    electricity_cost_post INTEGER,
    additional_incentive INTEGER,
    
    -- System details
    system_size DECIMAL(10,2) DEFAULT 0,
    panel_count INTEGER DEFAULT 0,
    panel_id INTEGER,
    inverter_id INTEGER,
    inverter_count INTEGER DEFAULT 1,
    battery_count INTEGER DEFAULT 0,
    
    -- Utility
    utility_id INTEGER,
    tariff_id INTEGER,
    
    -- Roof details
    roof_material INTEGER,
    surface_id INTEGER,
    
    -- Production
    annual_production DECIMAL(12,2) DEFAULT 0,
    
    -- Workflow states
    welcome_call_state INTEGER,
    financing_state INTEGER,
    utility_bill_state INTEGER,
    design_approved_state INTEGER,
    permitting_approved_state INTEGER,
    site_photos_state INTEGER,
    install_crew_state INTEGER,
    installation_state INTEGER,
    final_inspection_state INTEGER,
    pto_state INTEGER,
    
    -- Dates
    installation_date VARCHAR(20),
    date_ntp VARCHAR(20),
    date_installed VARCHAR(20)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_leads_external_lead_id ON leads(external_lead_id);
CREATE INDEX IF NOT EXISTS idx_leads_sync_status ON leads(sync_status);
CREATE INDEX IF NOT EXISTS idx_leads_company_id ON leads(company_id);
CREATE INDEX IF NOT EXISTS idx_leads_creator_id ON leads(creator_id);
CREATE INDEX IF NOT EXISTS idx_leads_state ON leads(state);
CREATE INDEX IF NOT EXISTS idx_leads_location ON leads(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_leads_panel_id ON leads(panel_id);
CREATE INDEX IF NOT EXISTS idx_leads_inverter_id ON leads(inverter_id);
CREATE INDEX IF NOT EXISTS idx_leads_utility_id ON leads(utility_id);
