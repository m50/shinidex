package dex

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/form"
	"github.com/m50/shinidex/pkg/web/middleware"
	"github.com/m50/shinidex/pkg/web/session"
)

func Router(e *echo.Echo) {
	g := e.Group("/dex", middleware.AuthnMiddleware, middleware.AuthzMiddleware)

	g.GET("", list)
	g.GET("/new", new)
	g.POST("", create)
	g.PUT("/:dex", update)
	g.DELETE("/:dex", delete)

	showRouter(g)
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
		Forms:         types.FormLocation(form.ParseInt(c.FormValue("forms"))),
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
	var form struct {
		id            string             `param:"dex"`
		name          string             `form:"name"`
		shiny         bool               `form:"shiny"`
		forms         types.FormLocation `form:"forms"`
		genderForms   types.FormLocation `form:"genderForms"`
		regionalForms types.FormLocation `form:"regionalForms"`
		gmaxForms     types.FormLocation `form:"gmaxForms"`
	}
	if err := c.Bind(form); err != nil {
		return err
	}

	db := c.(database.DBContext).DB().Pokedexes()
	dex, err := db.FindByID(form.id)
	if err != nil {
		return err
	}
	dex.Name = form.name
	if err = dex.UpdateConfig(types.PokedexConfig{
		Shiny:         form.shiny,
		Forms:         form.forms,
		GenderForms:   form.genderForms,
		RegionalForms: form.regionalForms,
		GMaxForms:     form.gmaxForms,
	}); err != nil {
		return err
	}

	if err = db.Update(dex); err != nil {
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/dex/%s", form.id))
}

func delete(c echo.Context) error {
	id := c.Param("dex")
	db := c.(database.DBContext).DB().Pokedexes()
	if err := db.Delete(id); err != nil {
		return err
	}

	return c.Redirect(http.StatusMovedPermanently, "/dex/")
}
