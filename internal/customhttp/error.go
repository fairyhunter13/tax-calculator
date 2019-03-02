package customhttp

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	ErrInvalidInput = echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
)
