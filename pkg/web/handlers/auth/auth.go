package auth

import "github.com/labstack/echo/v4"

func Router(e *echo.Echo) {
	group := e.Group("/auth")
	group.GET("/test", func (c echo.Context) error {
		return c.HTML(200, "Test")
	})
}