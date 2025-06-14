package static

import (
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo) {
	e.Any(GetScriptPath(), scriptRoute)
	e.Any(GetStylePath(), styleRoute)
}

func scriptRoute(c echo.Context) error {
	return c.File(scriptFSPath)
}

func styleRoute(c echo.Context) error {
	return c.File(styleFSPath)
}
