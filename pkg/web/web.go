package web

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/handlers/auth"
	"github.com/m50/shinidex/pkg/web/handlers/dex"
	"github.com/m50/shinidex/pkg/web/handlers/pokemon"
	smiddleware "github.com/m50/shinidex/pkg/web/middleware"
	"github.com/m50/shinidex/pkg/web/static"
	"github.com/spf13/viper"
)

type Context struct {
	echo.Context
	db *database.Database
}

func (c Context) DB() *database.Database {
	return c.db
}

func router(e *echo.Echo) {
	e.Static("/assets/imgs", "assets/imgs")
	e.Static("/icons", "icons")
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, e.Reverse("pokemon-list"))
	})
	static.Router(e)
	auth.Router(e)
	pokemon.Router(e)
	dex.Router(e)

	e.Any("*", func(c echo.Context) error {
		return views.RenderView(c, http.StatusNotFound, views.NotFound())
	})
}

func New(db *database.Database) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &Context{c, db}
			return next(cc)
		}
	})
	e.Use(smiddleware.LoggingHandler())
	e.Use(smiddleware.ErrorHandler)
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(viper.GetString("auth.key")))))
	e.Use(middleware.CORS())
	e.Use(smiddleware.HeaderMiddleware)

	router(e)

	return e
}
