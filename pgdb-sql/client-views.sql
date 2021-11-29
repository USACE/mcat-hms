-- Check views for upstream tables/materialized views before changing/using!
--DROP VIEW models.hms_project_summary;
CREATE OR REPLACE VIEW models.hms_project_summary AS 

SELECT  squery.col_1 AS "1. Project Title",
 		squery.col_2 AS "2. Description", 
 		squery.col_3 AS "3. Data Collection",
 		squery.col_4 AS "4. Source",
 		squery.s3_key AS s3_key
FROM 
	(SELECT
		t.title AS col_1,
		t.description AS col_2,
		i.title AS col_3,
		i."source" AS col_4,
		r.s3_key AS s3_key
			
	 FROM hms_definition_metadata t
	 JOIN models.model r ON r.model_inventory_id = t.model_inventory_id
	 JOIN inventory.collections i ON i.collection_id = t.collection 
	) squery;

-- DROP VIEW models.hms_met_files_view;
CREATE VIEW models.hms_met_files_view AS 
SELECT 
	hfm.title AS "1. Title",
	hfm.precip_type AS "2. Precipitation Type",
	hfm.met_file AS "3. Met File",
	hdm.s3_key AS "s3_key"
FROM 
	 models.hms_definition_metadata hdm, 
	 models.hms_forcing_metadata hfm 
WHERE hdm.model_inventory_id = hfm.model_inventory_id;

-- DROP VIEW models.hms_basin_files_view;
CREATE VIEW models.hms_basin_files_view AS 
SELECT 
	hgm.title AS "0. Basin Title",
	hgm.description AS "1. Description",
	'N/A' AS "2. Total Subbasin Area",
	hgm.num_subbasins AS "3. # of Subbasins",
	'N/A' AS "4. Average Subasin Area (sqmi)",
	hgm.loss_rate AS "5. Loss Method",
	hgm.transform_method AS "6. Transform Method",
	hgm.num_reaches AS "7. # of Reaches",
	hgm.num_junctions AS "8. # of Junctions",
	hgm.num_reservoirs AS "9. # of Reservoirs",
	hdm.s3_key AS "s3_key"
FROM 
	models.hms_definition_metadata hdm, 
	models.hms_geometry_metadata hgm 
WHERE hdm.model_inventory_id = hgm.model_inventory_id;

--DROP VIEW models.hms_control_files_view;
CREATE VIEW models.hms_control_files_view AS 
SELECT 
	hcm.title AS "1. Run Name",
	hfm.basin_model AS "2. Basin File",
	hfm.met_file AS "3. Precipitation File",
	hcm.description AS "4. Description",
	hdm.s3_key AS "s3_key"
FROM 
	 models.hms_definition_metadata hdm, 
	 models.hms_control_metadata hcm,
	 models.hms_forcing_metadata hfm
WHERE hdm.model_inventory_id = hcm.model_inventory_id
AND hdm.model_inventory_id = hfm.model_inventory_id;