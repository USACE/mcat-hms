package handlers

import (
	"net/http"

	hms "app/tools"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// GeospatialData godoc
// @Summary Extract geospatial data
// @Description Extract geospatial data from a HMS model given an s3 key
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {object} interface{}
// @Failure 500 {object} SimpleResponse
// @Router /geospatialdata [get]
func GeospatialData(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := hms.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error()})
		}
		data := hm.GeospatialData()

		return c.JSON(http.StatusOK, data)
	}
}
