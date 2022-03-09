-- Definition files
CREATE MATERIALIZED VIEW models.hms_definition_metadata
AS SELECT hms.model_inventory_id,
    c.collection_id AS collection,
    hms.model_metadata ->> 'Title'::text AS title,
    hms.model_metadata ->> 'Version'::text AS version,
    hms.model_metadata ->> 'Description'::text AS description,
    hms.s3_key
   FROM models.model hms
     LEFT JOIN inventory.collections c USING (collection_id)
  WHERE hms.type = 'HMS'::text
WITH DATA;

-- Control files 
CREATE MATERIALIZED VIEW models.hms_control_metadata
AS WITH control_metadata AS (
         SELECT m.model_inventory_id,
            d.key AS file,
            d.value AS meta
           FROM models.model m
             LEFT JOIN LATERAL json_each((m.model_metadata -> 'Metadata'::text) -> 'ControlMetadata'::text) d(key, value) ON true
          WHERE ((m.model_metadata -> 'Metadata'::text) ->> 'ControlMetadata'::text) IS NOT NULL
        ), model_metadata AS (
         SELECT model.model_metadata ->> 'ModelDirectory'::text AS dir
           FROM models.model
        )
 SELECT DISTINCT control_metadata.model_inventory_id,
    control_metadata.meta ->> 'Title'::text AS title,
    control_metadata.meta ->> 'Description'::text AS description,
    control_metadata.meta ->> 'Start Date'::text AS start_date,
    control_metadata.meta ->> 'Start Time'::text AS start_time,
    control_metadata.meta ->> 'End Date'::text AS end_date,
    control_metadata.meta ->> 'End Time'::text AS end_time,
    control_metadata.meta ->> 'Time Interval'::text AS time_interval,
    control_metadata.meta ->> 'Notes'::text AS notes,
    control_metadata.file AS control_file
   FROM control_metadata,
    model_metadata
WITH DATA;

-- Forcing files 
CREATE MATERIALIZED VIEW models.hms_forcing_metadata
AS WITH forcing_metadata AS (
         SELECT m.model_inventory_id,
            d.key AS file,
            d.value AS meta
           FROM models.model m
             JOIN LATERAL json_each((m.model_metadata -> 'Metadata'::text) -> 'ForcingMetadata'::text) d(key, value) ON true
          WHERE ((m.model_metadata -> 'Metadata'::text) ->> 'ForcingMetadata'::text) IS NOT NULL
        ), model_metadata AS (
         SELECT model.model_metadata ->> 'ModelDirectory'::text AS dir
           FROM models.model
        )
 SELECT DISTINCT forcing_metadata.model_inventory_id,
    forcing_metadata.meta ->> 'Title'::text AS title,
    forcing_metadata.meta ->> 'Description'::text AS description,
    forcing_metadata.meta ->> 'Unit System'::text AS units,
    (forcing_metadata.meta -> 'Use Basin Model'::text) ->> 0 AS basin_model,
    forcing_metadata.meta ->> 'Precipitation'::text AS precip_type,
    json_array_length(
        CASE
            WHEN ((forcing_metadata.meta -> 'Subbasin'::text)::text) = 'null'::text THEN '[]'::json
            ELSE forcing_metadata.meta -> 'Subbasin'::text
        END) AS num_subbasins,
    forcing_metadata.meta ->> 'Notes'::text AS notes,
    forcing_metadata.file AS met_file
   FROM forcing_metadata,
    model_metadata
WITH DATA;

-- Geometry files 
CREATE MATERIALIZED VIEW models.hms_geometry_metadata
AS WITH geometry_metadata AS (
         SELECT m.model_inventory_id,
            d.key AS file,
            d.value AS meta
           FROM models.model m
             JOIN LATERAL json_each((m.model_metadata -> 'Metadata'::text) -> 'GeometryMetadata'::text) d(key, value) ON true
          WHERE ((m.model_metadata -> 'Metadata'::text) ->> 'GeometryMetadata'::text) IS NOT NULL
        ), model_metadata AS (
         SELECT model.model_metadata ->> 'ModelDirectory'::text AS dir
           FROM models.model
        )
 SELECT DISTINCT geometry_metadata.model_inventory_id,
    geometry_metadata.meta ->> 'Title'::text AS title,
    geometry_metadata.meta ->> 'Description'::text AS description,
    geometry_metadata.meta ->> 'Unit System'::text AS units,
    geometry_metadata.meta ->> 'LossRate'::text AS loss_rate,
    geometry_metadata.meta ->> 'Transform'::text AS transform_method,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Subbasin'::text), 0) AS num_subbasins,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Reach'::text), 0) AS num_reaches,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Junction'::text), 0) AS num_junctions,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Source'::text), 0) AS num_sources,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Sink'::text), 0) AS num_sinks,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Reservoir'::text), 0) AS num_reservoirs,
    COALESCE(json_array_length((geometry_metadata.meta -> 'Features'::text) -> 'Diversion'::text), 0) AS num_diversions,
    geometry_metadata.meta ->> 'Notes'::text AS notes,
    geometry_metadata.file AS basin_file
   FROM geometry_metadata,
    model_metadata
WITH DATA;


-- ALTER MATERIALIZED VIEW models.hms_definition_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_control_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_forcing_metadata OWNER TO readwrite;
-- ALTER MATERIALIZED VIEW models.hms_geometry_metadata OWNER TO readwrite;