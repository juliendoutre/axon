CREATE TABLE IF NOT EXISTS axon.observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_type TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    attributes JSONB NOT NULL,
    observer_claims JSONB NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS observations_asset_type_index
ON axon.observations USING gin (to_tsvector('english', asset_type));

CREATE INDEX IF NOT EXISTS observations_asset_id_index
ON axon.observations USING gin (to_tsvector('english', asset_id));

CREATE INDEX IF NOT EXISTS observations_attributes_index
ON axon.observations USING gin (attributes jsonb_path_ops);

CREATE INDEX IF NOT EXISTS observations_observer_claims_index
ON axon.observations USING gin (observer_claims jsonb_path_ops);

CREATE INDEX IF NOT EXISTS observations_timestamp_index
ON axon.observations USING btree (timestamp);
