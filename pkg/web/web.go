package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/pkg/web/handlers/auth"
	"github.com/m50/shinidex/pkg/web/handlers/pokemon"
)

func router(e *echo.Echo) {
	e.Static("/assets", "assets")
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, e.Reverse("pokemon-list"))
	})
	auth.Router(e)
	pokemon.Router(e)
}

func New() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} \x1b[34mRQST\x1b[0m ${method} http://${host}${uri} : ${status} ${error}\n",
	}))

	router(e)

	return e
}
