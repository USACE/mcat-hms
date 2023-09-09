package handlers

import (
	"fmt"
	"net/http"

	"github.com/Dewberry/mcat-hms/config"

	"github.com/USACE/filestore"
	"github.com/labstack/echo/v4"
)

type SimpleResponse struct {
	Status  int
	Message string
}

// Ping godoc
// @Summary Status Check
// @Description Check which services are operational
// @Tags Health Check
// @Accept  json
// @Produce  json
// @Success 200 {object} SimpleResponse
// @Router /ping [get]
func Ping(ac *config.APIConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := ac.DB.Ping()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unavailable"})
		}
		switch (*ac.FileStore).(type) {
		case *filestore.BlockFS:
			fmt.Println("File is local")
			return c.JSON(http.StatusOK, map[string]string{"status": "available"})

		case *filestore.S3FS:
			s3FS := (*ac.FileStore).(*filestore.S3FS)
			err := s3FS.Ping()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unavailable"})
			}
			fmt.Println("File is on S3")
			return c.JSON(http.StatusOK, map[string]string{"status": "available"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unavailable"})
	}
}
