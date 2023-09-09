package config

import (
	"fmt"
	"os"

	"github.com/USACE/filestore"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type APIConfig struct {
	Host      string
	Port      int
	FileStore *filestore.FileStore
	DB        *sqlx.DB
}

// Address tells the application where to run the api out of
func (app *APIConfig) Address() string {
	return fmt.Sprintf("%s:%d", app.Host, app.Port)
}

// Init initializes the API's configuration
func Init() *APIConfig {
	config := new(APIConfig)
	config.Host = "" // 0.0.0.0
	config.Port = 5900
	config.FileStore = FileStoreInit(false)
	config.DB = DBInit()
	return config
}

// FileStoreInit initializes the filestore object
func FileStoreInit(local bool) *filestore.FileStore {

	var fs filestore.FileStore
	var err error
	switch local {
	case true:
		fs, err = filestore.NewFileStore(filestore.BlockFSConfig{})
		if err != nil {
			panic(err)
		}
	case false:
		config := filestore.S3FSConfig{
			S3Id:     os.Getenv("AWS_ACCESS_KEY_ID"),
			S3Key:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			S3Region: os.Getenv("AWS_DEFAULT_REGION"),
			S3Bucket: os.Getenv("S3_BUCKET"),
		}

		fs, err = filestore.NewFileStore(config)
		if err != nil {
			panic(err)
		}
	}
	return &fs
}

func DBInit() *sqlx.DB {

	creds := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))

	return sqlx.MustOpen("pgx", creds)
}
