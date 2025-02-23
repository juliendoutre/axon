CREATE TABLE IF NOT EXISTS axon.extracted_assets (
    observation_id UUID REFERENCES axon.observations (id),
    attributes_path TEXT NOT NULL,
    asset_type TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS extracted_assets_attributes_path_index
ON axon.extracted_assets USING gin (to_tsvector('english', attributes_path));

CREATE INDEX IF NOT EXISTS extracted_assets_asset_type_index
ON axon.extracted_assets USING gin (to_tsvector('english', asset_type));

CREATE INDEX IF NOT EXISTS extracted_assets_asset_id_index
ON axon.extracted_assets USING gin (to_tsvector('english', asset_id));

CREATE INDEX IF NOT EXISTS extracted_assets_timestamp_index
ON axon.extracted_assets USING btree (timestamp);
