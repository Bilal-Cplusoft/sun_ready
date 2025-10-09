-- Migration: Add LightFusion 3D model fields to leads table
-- Purpose: Track 3D models created via LightFusion API for leads

-- Add columns for LightFusion 3D project tracking
ALTER TABLE leads 
ADD COLUMN IF NOT EXISTS lightfusion_3d_project_id INTEGER,
ADD COLUMN IF NOT EXISTS lightfusion_3d_house_id INTEGER,
ADD COLUMN IF NOT EXISTS model_3d_status VARCHAR(50),
ADD COLUMN IF NOT EXISTS model_3d_created_at TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS model_3d_completed_at TIMESTAMPTZ;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_leads_lightfusion_3d_project_id ON leads(lightfusion_3d_project_id);
CREATE INDEX IF NOT EXISTS idx_leads_lightfusion_3d_house_id ON leads(lightfusion_3d_house_id);
CREATE INDEX IF NOT EXISTS idx_leads_model_3d_status ON leads(model_3d_status);

-- Add comments for documentation
COMMENT ON COLUMN leads.lightfusion_3d_project_id IS 'LightFusion API project ID for 3D model';
COMMENT ON COLUMN leads.lightfusion_3d_house_id IS 'LightFusion API house ID for 3D model';
COMMENT ON COLUMN leads.model_3d_status IS '3D model status from LightFusion: pending, processing, completed, failed';
COMMENT ON COLUMN leads.model_3d_created_at IS 'Timestamp when 3D model creation was initiated';
COMMENT ON COLUMN leads.model_3d_completed_at IS 'Timestamp when 3D model was completed';