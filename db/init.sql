-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    code VARCHAR(255) NOT NULL DEFAULT '',
    slug VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    logo_path TEXT,
    admin_id INTEGER,
    sales_commission_min DECIMAL(10,4),
    sales_commission_max DECIMAL(10,4),
    sales_commission_default DECIMAL(10,4),
    baseline DECIMAL(10,2),
    baseline_adder DECIMAL(10,2),
    baseline_adder_pct_sales_comms INTEGER,
    contract_tag VARCHAR(100),
    referred_by_user_id INTEGER,
    credits INTEGER,
    custom_commissions BOOLEAN NOT NULL DEFAULT false,
    pricing_mode INTEGER NOT NULL DEFAULT 0
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    firstname VARCHAR(200),
    lastname VARCHAR(200),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(100),
    type SMALLINT NOT NULL,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE NO ACTION,
    creator_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    picture_path TEXT,
    disabled BOOLEAN NOT NULL DEFAULT false,
    is_manager BOOLEAN NOT NULL DEFAULT false
);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    address TEXT
);

-- Create leads table
CREATE TABLE IF NOT EXISTS leads (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- External sync fields
    external_lead_id INTEGER UNIQUE,
    sync_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    last_synced_at TIMESTAMPTZ,
    
    -- Core lead information
    state INTEGER NOT NULL DEFAULT 0,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    creator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Location data
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    address TEXT,
    
    -- Source and metadata
    source INTEGER NOT NULL DEFAULT 0,
    promo_code VARCHAR(100),
    is_2d BOOLEAN NOT NULL DEFAULT false,
    
    -- Energy consumption details
    kwh_usage DECIMAL(10, 2),
    kwh_per_kw_manual INTEGER,
    
    -- Financial information
    electricity_cost_pre INTEGER,
    electricity_cost_post INTEGER,
    additional_incentive INTEGER,
    
    -- System specifications
    system_size DECIMAL(10, 2),
    panel_count INTEGER,
    panel_id INTEGER,
    inverter_id INTEGER,
    inverter_count INTEGER DEFAULT 1,
    battery_count INTEGER DEFAULT 0,
    
    -- Utility information
    utility_id INTEGER,
    tariff_id INTEGER,
    
    -- Roof details
    roof_material INTEGER,
    surface_id INTEGER,
    
    -- Production metrics
    annual_production DECIMAL(12, 2),
    
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
    
    -- Important dates
    installation_date DATE,
    date_ntp DATE,
    date_installed DATE,
    
    -- LightFusion 3D Project Integration
    lightfusion_3d_project_id INTEGER,
    lightfusion_3d_house_id INTEGER,
    model_3d_status VARCHAR(50),
    model_3d_created_at TIMESTAMPTZ,
    model_3d_completed_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT chk_latitude CHECK (latitude >= -90 AND latitude <= 90),
    CONSTRAINT chk_longitude CHECK (longitude >= -180 AND longitude <= 180),
    CONSTRAINT chk_sync_status CHECK (sync_status IN ('pending', 'synced', 'failed', 'syncing'))
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_projects_company_id ON projects(company_id);
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_companies_slug ON companies(slug);
CREATE INDEX IF NOT EXISTS idx_leads_company_id ON leads(company_id);
CREATE INDEX IF NOT EXISTS idx_leads_creator_id ON leads(creator_id);
CREATE INDEX IF NOT EXISTS idx_leads_external_lead_id ON leads(external_lead_id);
CREATE INDEX IF NOT EXISTS idx_leads_sync_status ON leads(sync_status);
CREATE INDEX IF NOT EXISTS idx_leads_lightfusion_3d_project_id ON leads(lightfusion_3d_project_id);
CREATE INDEX IF NOT EXISTS idx_leads_model_3d_status ON leads(model_3d_status);

-- Insert a default company
INSERT INTO companies (name, display_name, description, code, slug, is_active)
VALUES ('Default Company', 'Default Company', 'Default company for testing', 'DEFAULT', 'default', true)
ON CONFLICT (slug) DO NOTHING;
