package pgdb

import (
	"app/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Upsert ...
func Upsert(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")
		if definitionFile == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": http.StatusBadRequest, "message": "Missing query parameter: `definition_file`"})
		}

		optimizedbStr := c.QueryParam("optimizedb")
		if optimizedbStr == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": http.StatusBadRequest, "message": "Missing query parameter: `optimizedb`"})
		}

		optimizedb, err := strconv.ParseBool(optimizedbStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": http.StatusBadRequest, "message": "Query parameter `optimizedb` is not a valid bool: " + err.Error()})
		}

		err = UpsertToDB(definitionFile, ac)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": http.StatusInternalServerError, "message": err.Error()})
		}

		if optimizedb {
			for _, query := range VacuumQuery {
				ac.DB.MustExec(query)
			}

			for _, query := range RefreshViewsQuery {
				ac.DB.MustExec(query)
			}
		}

		return c.JSON(http.StatusOK, "Successfully added HMS model metadata to the database: "+definitionFile)
	}
}

// Upsert ...
func UpsertHMSGeometry(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}

// Upsert ...
func RefreshHMSViews(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}

// Upsert ...
func VacuumHMSViews(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}
