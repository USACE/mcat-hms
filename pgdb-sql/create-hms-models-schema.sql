CREATE SCHEMA IF NOT EXISTS models;

/*---------------------------------------------------------------------------*/
-- Create models.column_reference table
/*---------------------------------------------------------------------------*/
CREATE TABLE IF NOT EXISTS models.column_reference(
    model_type TEXT PRIMARY KEY,
    map JSON NOT NULL
);


/*---------------------------------------------------------------------------*/
-- Create models.model table
/*---------------------------------------------------------------------------*/
CREATE TABLE IF NOT EXISTS models.model (
    model_inventory_id BIGINT PRIMARY KEY,
    collection_id BIGINT,
    name TEXT,
    type TEXT,
    s3_key TEXT UNIQUE NOT NULL,
    model_metadata JSON NOT NULL,
    etl_metadata JSON NOT NULL,
    CONSTRAINT model_collection_id_fk FOREIGN KEY (collection_id) REFERENCES inventory.collections (collection_id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT model_type_check CHECK (
        type = 'RAS' OR
        type = 'HMS' OR
        type = 'OTHER')
);

-- Create indexes on foreign keys
CREATE INDEX IF NOT EXISTS collection_id_idx ON models.model (collection_id);
CREATE INDEX IF NOT EXISTS model_type_idx ON models.model (type);
