package handlers

import (
	"net/http"

	"github.com/Dewberry/mcat-hms/tools"
	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// ModelVersion godoc
// @Summary Extract the HMS model version
// @Description Extract the HMS model version given an s3 key
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {string} string "4.0"
// @Failure 500 {object} SimpleResponse
// @Router /modelversion [get]
func ModelVersion(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := tools.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error()})
		}
		vers := hm.ModelVersion()

		return c.JSON(http.StatusOK, vers)
	}
}
