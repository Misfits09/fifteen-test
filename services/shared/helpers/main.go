package helpers

import (
	"fifteen/shared/structs"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendErrorResponse(c echo.Context, err error, status int) error {
	if err == nil {
		return c.JSON(status, &structs.ErrorResponse{
			Success: false,
			Error:   http.StatusText(status),
		})
	} else {
		return c.JSON(status, &structs.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}
}

func LogIfIsError(err error) bool {
	if err != nil {
		log.Print(err)
		return true
	}
	return false
}
