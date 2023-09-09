package handlers

import (
	"net/http"

	"github.com/Dewberry/mcat-hms/tools"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// IsGeospatial godoc
// @Summary Check if the HMS model has geospatial information
// @Description Check if the HMS model has geospatial information
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {object} bool
// @Router /isgeospatial [get]
func IsGeospatial(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := tools.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error()})
		}
		isIt := hm.IsGeospatial()

		return c.JSON(http.StatusOK, isIt)
	}
}
