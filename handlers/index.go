package handlers

import (
	"net/http"

	"github.com/Dewberry/mcat-hms/tools"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// Index godoc
// @Summary Index a HMS model
// @Description Extract metadata from a HMS model given an s3 key
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {object} hms.Model
// @Failure 500 {object} SimpleResponse
// @Router /index [get]
func Index(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := tools.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, SimpleResponse{http.StatusInternalServerError, err.Error()})
		}
		mod := hm.Index()

		return c.JSON(http.StatusOK, mod)
	}
}
