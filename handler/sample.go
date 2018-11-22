package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func SampleString(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func SampleJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"hello": "world"})
}

func SampleHTML(c echo.Context) error {
	return c.Render(http.StatusOK, "sample/index", map[string]interface{}{"Name": "guest"})
}
