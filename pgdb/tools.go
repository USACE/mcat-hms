package pgdb

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/Dewberry/mcat-hms/config"
	"github.com/Dewberry/mcat-hms/tools"

	"github.com/jmoiron/sqlx"
)

// func getETLMetadata(hm *tools.HmsModel) ([]byte, error) {
// 	metadataFiles := make([]string, 0)

// 	prefix := hm.ModelDirectory + "/"

// 	files, err := hm.FileStore.GetDir(prefix, false)
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	for _, file := range *files {
// 		filebase := file.Name
// 		if strings.Contains(filebase, ".") {
// 			fileparts := strings.Split(filebase, ".")
// 			nparts := len(fileparts)
// 			ending := strings.Join([]string{fileparts[nparts-2], fileparts[nparts-1]}, ".")
// 			if ending == "metadata.json" {
// 				metadataFiles = append(metadataFiles, tools.BuildFilePath(hm.ModelDirectory, file.Name))
// 			}
// 		}
// 	}

// 	if len(metadataFiles) == 0 {
// 		return []byte{}, errors.New("file not found: model etl metadata json")
// 	} else if len(metadataFiles) > 1 {
// 		return []byte{}, errors.New("multiple files found: model etl metadata json")
// 	}

// 	jsonFile, err := hm.FileStore.GetObject(metadataFiles[0])
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	defer jsonFile.Close()

// 	return ioutil.ReadAll(jsonFile)
// }

// `SELECT c.collection_id FROM models.collections c
// INNER JOIN inv.sources s ON c.source_id = s.source_id AND s.source = $1
// WHERE c.collection = $2;`
func getCollectionID(tx *sqlx.Tx, definitionFile string) (int, error) {

	var collectionID int

	if err := tx.Get(&collectionID, collectionIDQuery, definitionFile); err != nil {
		return 0, err

	}

	return collectionID, nil
}

// `INSERT INTO models.model (name, type, collection_id, s3_key, model_metadata, etl_metadata) VALUES ($1, $2, $3, $4, $5, $6)
// ON CONFLICT (s3_key)
// DO UPDATE SET name = $1, type = $2, collection_id = $3, s3_key = $4, model_metadata = $5, etl_metadata = $6 RETURNING model_inventory_id;`
func upsertModels(tx *sqlx.Tx, collectionID int, modelName string, definitionFile string, hmMarshal []byte, etlMetadata []byte) error {
	modQuery := MakeUpsertQuery(QueryConfig{"models", "model", "model_inventory_id", []string{"s3_key"}, []string{"name", "type", "collection_id",
		"s3_key", "model_metadata", "etl_metadata"}}, false)

	tx.MustExec(modQuery, modelName, "HMS", collectionID, definitionFile, hmMarshal, etlMetadata)

	return nil
}

// Upsert extracts metadata from a HEC-HMS model and writes this data to a database using upsert statements. These
// upsert statements insert a row where it does not exist, or updates the row with new values when it does. A database
// transaction is used so that in the case of an error, no model data is written to the database, i.e. the transation is
// aborted, which prevents incomplete tables for a given model. Furthermore, tx.Get is used so that the primary key from
// one table is returned and can be used as a foreign key in subsequent tables.
func UpsertToDB(definitionFile string, ac *config.APIConfig) (err error) {
	defFileName := filepath.Base(definitionFile)
	modelName := strings.TrimSuffix(defFileName, filepath.Ext(defFileName))

	hm, err := tools.NewHmsModel(definitionFile, *ac.FileStore)
	if err != nil {
		return
	}

	hmMarshal, err := json.Marshal(hm)
	if err != nil {
		return
	}

	// etlMetadata, err := getETLMetadata(hm)
	// if err != nil {
	// 	return
	// }

	etlMetadata := map[string]string{"model_name": hm.Title, "source_path": definitionFile, "destination_path": "", "projection_source_path": ""}
	etlMarshal, err := json.Marshal(etlMetadata)
	if err != nil {
		return
	}

	tx := ac.DB.MustBegin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	collectionID, err := getCollectionID(tx, definitionFile)
	if err != nil {
		return
	}

	err = upsertModels(tx, collectionID, modelName, definitionFile, hmMarshal, etlMarshal)
	if err != nil {
		return
	}

	return
}
