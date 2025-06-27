package error

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func InvalidRequestBody(msg any) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body: %s", msg))
}

func InternalServerError(msg any) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, msg)
}

func Unauthorized(msg any) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, msg)
}
