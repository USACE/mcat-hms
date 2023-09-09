package pgdb

import (
	"net/http"

	"github.com/Dewberry/mcat-hms/config"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// Upsert ...
func UpsertHMSModel(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")
		if definitionFile == "" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"status":  http.StatusBadRequest,
					"message": "Missing query parameter: `definition_file`",
				},
			)
		}

		err := UpsertToDB(definitionFile, ac)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				},
			)
		}

		return c.JSON(http.StatusOK, "Successfully added HMS model metadata to the database: "+definitionFile)
	}
}

// // UpsertHMSGeometry ...
// func UpsertHMSGeometry(ac *config.APIConfig) echo.HandlerFunc {
// 	return func(c echo.Context) error { return nil }
// }

// VacuumHMSViews ...
func VacuumHMSViews(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		for _, query := range vacuumQuery {
			_, err := db.Exec(query)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
		}

		return c.JSON(http.StatusOK, "HMS tables vacuumed successfully.")
	}
}

// RefreshHMSViews ...
func RefreshHMSViews(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		for _, query := range refreshViewsQuery {
			_, err := db.Exec(query)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
		}

		return c.JSON(http.StatusOK, "HMS materialized views refreshed successfully.")
	}
}
