package dex

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/form"
	"github.com/m50/shinidex/pkg/web/session"
)

func Router(e *echo.Echo) {
	g := e.Group("/dex", views.AuthnMiddleware())

	g.GET("", list)
	g.GET("/:dex", show)
	g.GET("/new", new)
	g.POST("", create)
	g.PUT("/:dex", update)
	g.DELETE("/:dex", delete)
}

func list(c echo.Context) error {
	user, err := session.GetAuthedUser(c)
	if err != nil {
		return err
	}
	db := c.(database.DBContext).DB()
	dexes, err := db.Pokedexes().FindByOwnerID(user.ID)
	if err != nil {
		return err
	}
	return views.RenderView(c, http.StatusOK, List(dexes))
}

func show(c echo.Context) error {
	return nil
}

func new(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, New())
}

func create(c echo.Context) error {
	user, err := session.GetAuthedUser(c)
	if err != nil {
		return err
	}
	name := c.FormValue("name")
	cfg := types.PokedexConfig{
		Shiny:         form.ParseBool(c.FormValue("shiny")),
		GenderForms:   types.FormLocation(form.ParseInt(c.FormValue("gender"))),
		RegionalForms: types.FormLocation(form.ParseInt(c.FormValue("regional"))),
		GMaxForms:     types.FormLocation(form.ParseInt(c.FormValue("gmax"))),
	}
	dex, err := types.NewPokedex(user.ID, name, cfg)
	if err != nil {
		return err
	}
	db := c.(database.DBContext).DB().Pokedexes()
	id, err := db.Insert(dex)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/dex/%s", id))
}

func update(c echo.Context) error {
	return nil
}

func delete(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/dex/")
}
