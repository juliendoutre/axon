CREATE TABLE IF NOT EXISTS axon.extracted_assets (
    observation_id UUID REFERENCES axon.observations (id),
    attributes_path TEXT NOT NULL,
    asset_type TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
