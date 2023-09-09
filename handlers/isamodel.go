package handlers

import (
	"net/http"

	"github.com/Dewberry/mcat-hms/tools"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

// IsAModel godoc
// @Summary Check if the given key is a HMS model
// @Description Check if the given key is a HMS model
// @Tags MCAT
// @Accept json
// @Produce json
// @Param definition_file query string true "/models/hms/Truckee_River/Truckee_River.hms"
// @Success 200 {object} bool
// @Router /isamodel [get]
func IsAModel(fs *filestore.FileStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		definitionFile := c.QueryParam("definition_file")

		hm, err := tools.NewHmsModel(definitionFile, *fs)
		if err != nil {
			return c.JSON(http.StatusOK, false)
		}
		isIt := hm.IsAModel()

		return c.JSON(http.StatusOK, isIt)
	}
}
