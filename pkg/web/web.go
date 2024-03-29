package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/web/handlers/auth"
	"github.com/m50/shinidex/pkg/web/handlers/pokemon"
)

type Context struct {
	echo.Context
	db *database.Database
}

func (c Context) DB() *database.Database {
	return c.db
}

func router(e *echo.Echo) {
	e.Static("/assets", "assets")
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, e.Reverse("pokemon-list"))
	})
	auth.Router(e)
	pokemon.Router(e)
}

func logger() *log.Logger {
	logger := log.New("shinidex")
	logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line}")
	logger.SetLevel(log.DEBUG)
	return logger
}

func New(db *database.Database) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Logger = logger()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{c, db}
			return next(cc)
		}
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} \x1b[34mRQST\x1b[0m ${method} http://${host}${uri} : ${status} ${error}\n",
	}))

	router(e)

	return e
}
