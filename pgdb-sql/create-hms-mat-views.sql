-- Definition files
CREATE MATERIALIZED VIEW models.hms_definition_metadata AS
SELECT
    hms.model_inventory_id,
    c.collection_id AS collection,
    (hms.model_metadata ->> 'Title' ) AS title,
    (hms.model_metadata ->> 'Version' ) AS version,
    (hms.model_metadata ->> 'Description' ) AS description,
    hms.s3_key AS s3_key
FROM models.model AS hms
LEFT JOIN inventory.collections AS c USING (collection_id)
WITH DATA;

-- Control files 
CREATE MATERIALIZED VIEW models.hms_control_metadata AS
WITH control_metadata as (
    SELECT hms.model_inventory_id, d.key as file, d.value as meta
    FROM models.model hms
    JOIN json_each((hms.model_metadata -> 'Metadata' -> 'ControlMetadata')) d ON true
    WHERE (hms.model_metadata -> 'Metadata' ->> 'ControlMetadata') IS NOT NULL
), model_metadata as (SELECT (model_metadata ->> 'ModelDirectory') as dir FROM models.model)
SELECT 
    control_metadata.model_inventory_id, 
    (control_metadata.meta ->> 'Title') as title, 
    (control_metadata.meta ->> 'Description') as description,
    (control_metadata.meta ->> 'Start Date') as start_date,
    (control_metadata.meta ->> 'Start Time') as start_time,
    (control_metadata.meta ->> 'End Date') as end_date,
    (control_metadata.meta ->> 'End Time') as end_time,
    (control_metadata.meta ->> 'Time Interval') as time_interval,
    (control_metadata.meta ->> 'Notes') as notes,
    model_metadata.dir || '/' || control_metadata.file as s3_key
FROM control_metadata, model_metadata
WITH DATA;

-- Forcing files 
CREATE MATERIALIZED VIEW models.hms_forcing_metadata AS
WITH forcing_metadata as (
    SELECT hms.model_inventory_id, d.key as file, d.value as meta 
    FROM models.model hms
    JOIN json_each((hms.model_metadata -> 'Metadata' -> 'ForcingMetadata')) d ON true
    WHERE (hms.model_metadata -> 'Metadata' ->> 'ForcingMetadata') IS NOT NULL
), model_metadata as (SELECT (model_metadata ->> 'ModelDirectory') as dir FROM models.model)
SELECT 
    forcing_metadata.model_inventory_id, 
    (forcing_metadata.meta ->> 'Title') as title, 
    (forcing_metadata.meta ->> 'Description') as description,
    (forcing_metadata.meta ->> 'Unit System') as units,
    (forcing_metadata.meta ->> 'Use Basin Model') as basin_model,
    json_array_length(CASE WHEN (forcing_metadata.meta -> 'Subbasin')::text = 'null' THEN '[]'::json ELSE (forcing_metadata.meta -> 'Subbasin') END) as num_subbasins,
    (forcing_metadata.meta ->> 'Notes') as notes,
    model_metadata.dir || '/' || forcing_metadata.file as s3_key
FROM forcing_metadata, model_metadata
WITH DATA;

-- Geometry files 
CREATE MATERIALIZED VIEW models.hms_geometry_metadata AS
WITH geometry_metadata as (
    SELECT hms.model_inventory_id, d.key as file, d.value as meta
    FROM models.model hms
    JOIN json_each((hms.model_metadata -> 'Metadata' -> 'GeometryMetadata')) d ON true
    WHERE (hms.model_metadata -> 'Metadata' ->> 'GeometryMetadata') IS NOT NULL
), model_metadata as (SELECT (model_metadata ->> 'ModelDirectory') as dir FROM models.model)
SELECT 
    geometry_metadata.model_inventory_id, 
    (geometry_metadata.meta ->> 'Title') as title, 
    (geometry_metadata.meta ->> 'Description') as description,
    (geometry_metadata.meta ->> 'Unit System') as units,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Subbasin')), 0) as num_subbasins,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Reach')), 0) as num_reaches,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Junction')), 0) as num_junctions,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Source')), 0) as num_sources,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Sink')), 0) as num_sinks,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Reservoir')), 0) as num_reservoirs,
    coalesce((json_array_length(geometry_metadata.meta -> 'Features' -> 'Diversion')), 0) as num_diversions,
    (geometry_metadata.meta ->> 'Notes') as notes,
    model_metadata.dir || '/' || geometry_metadata.file as s3_key
FROM geometry_metadata, model_metadata
WITH DATA;


-- ALTER MATERIALIZED VIEW models.hms_definition_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_control_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_forcing_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_geometry_metadata OWNER TO readwrite;