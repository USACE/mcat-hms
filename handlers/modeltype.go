package handlers

import (
	"net/http"

	hms "github.com/Dewberry/mcat-hms/tools"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// ModelType godoc
// @Summary Extract the model type
// @Description Extract the model type given an s3 key
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {string} string "HMS"
// @Failure 500 {object} SimpleResponse
// @Router /modeltype [get]
func ModelType(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := hms.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error()})
		}
		typ := hm.ModelType()

		return c.JSON(http.StatusOK, typ)
	}
}
