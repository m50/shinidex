package pokemon

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/views"
)

var clicked int

func Router(e *echo.Echo) {
	clicked = 0
	group := e.Group("/pokemon")

	group.GET("/", list).Name = "pokemon-list"
	group.GET("/click", click)
}

func list(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, Main())
}

func click(c echo.Context) error {
	clicked++
	return views.RenderView(c, http.StatusOK, Click(clicked))
}