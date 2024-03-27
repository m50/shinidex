package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/frontend/views"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} \x1b[34mRQST\x1b[0m ${method} http://${host}${uri} : ${status} ${error}\n",
	}))

	e.Static("/assets", "assets")
	e.GET("/", func(c echo.Context) error {
		return views.RenderView(c, http.StatusOK, views.Home())
	})
	clicked := 0
	e.GET("/click", func(c echo.Context) error {
		clicked++
		return views.RenderView(c, http.StatusOK, views.Click(clicked))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
