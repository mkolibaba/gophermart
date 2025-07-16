package error

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InvalidRequestBody(msg any) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", msg))
}

func InternalServerError(err error) error {
	return err
}

func Unauthorized(msg any) error {
	return echo.NewHTTPError(http.StatusUnauthorized, msg)
}
