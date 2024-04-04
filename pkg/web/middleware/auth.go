package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/errors"
	"github.com/m50/shinidex/pkg/web/session"
)

func AuthnMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !session.IsLoggedIn(c) {
			return views.RenderView(c, http.StatusUnauthorized, views.Unauthorized())
		}
		return next(c)
	}
}

func AuthzMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !session.IsLoggedIn(c) {
			return next(c)
		}

		dexID := c.Param("dex")
		if dexID == "" {
			return next(c)
		}
		db := c.(database.DBContext).DB()
		dex, err := db.Pokedexes().FindByID(dexID)
		if err != nil {
			return err
		}

		user, err := session.GetAuthedUser(c)
		if err != nil {
			return err
		}

		if dex.OwnerID != user.ID {
			return errors.NewForbiddenError("unable to access this pokedex")
		}

		return next(c)
	}
}

func HeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return err
		}
		if c.Request().Header.Get("hx-request") != "true" {
			return nil
		}
		user, _ := session.GetAuthedUser(c)
		return views.AddView(c, views.Header(user))
	}
}
