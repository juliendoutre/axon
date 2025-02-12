CREATE TABLE IF NOT EXISTS axon.observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_type TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    attributes JSONB NOT NULL,
    observer_claims JSONB NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
