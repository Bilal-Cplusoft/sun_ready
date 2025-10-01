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

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_projects_company_id ON projects(company_id);
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_companies_slug ON companies(slug);

-- Insert a default company
INSERT INTO companies (name, display_name, description, code, slug, is_active)
VALUES ('Default Company', 'Default Company', 'Default company for testing', 'DEFAULT', 'default', true)
ON CONFLICT (slug) DO NOTHING;
