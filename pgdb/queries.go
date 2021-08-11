package pgdb

import (
	"fmt"
	"os"
	"strings"
)

var collectionIDQuery = fmt.Sprintf(`
SELECT collection_id 
FROM inventory.collections 
WHERE 's3://%s/'`, os.Getenv("S3_BUCKET")+` || $1 LIKE s3_prefix || '%';`)

// VacuumQuery ...
var vacuumQuery []string = []string{"VACUUM ANALYZE models.model;"}

// RefreshViewsQuery ...
var refreshViewsQuery []string = []string{"REFRESH MATERIALIZED VIEW models.hms_definition_metadata;",
	"REFRESH MATERIALIZED VIEW models.hms_control_metadata;",
	"REFRESH MATERIALIZED VIEW models.hms_forcing_metadata;",
	"REFRESH MATERIALIZED VIEW models.hms_geometry_metadata;"}

// QueryConfig ...
type QueryConfig struct {
	schemaName string
	tableName  string
	primaryKey string
	uniqueKeys []string
	keys       []string
}

// MakeUpsertQuery builds the upsert query using the QueryConfig struct. The following is an example of the upsert query
// used for writing data to the models table:
// 'INSERT INTO models.hms (name, data_group, type, collection_id, s3_key, model_metadata, etl_metadata) VALUES ($1, $2, $3, $4, $5, $6, $7)
// ON CONFLICT (s3_key)
// DO UPDATE SET name = $1, data_group = $2, type = $3, collection_id = $4, s3_key = $5, model_metadata = $6, etl_metadata = $7;'
func MakeUpsertQuery(qc QueryConfig, returnPK bool) string {
	cols := qc.keys[0]
	vals := "$1"
	validx := 2
	colVals := fmt.Sprintf("%s = %s", cols, vals)

	for _, key := range qc.keys[1:] {
		cols += fmt.Sprintf(", %s", key)
		vals += fmt.Sprintf(", $%d", validx)
		colVals += fmt.Sprintf(", %s = $%d", key, validx)
		validx++
	}
	conflictCols := strings.Join(qc.uniqueKeys, ",")

	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s;", qc.schemaName, qc.tableName, cols, vals, conflictCols, colVals)
	if returnPK {
		return strings.Replace(query, ";", fmt.Sprintf(" RETURNING %s;", qc.primaryKey), 1)

	}
	return query
}
