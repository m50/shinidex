package dex

import (
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
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
	g.GET("/:dex/edit", edit)
	g.POST("", create)
	g.PUT("/:dex", update)
	g.DELETE("/:dex", delete)

	showRouter(g)
}

func list(c echo.Context) error {
	ctx := context.FromEcho(c)
	user, err := session.GetAuthedUser(c)
	if err != nil {
		return err
	}
	db := c.(database.DBContext).DB()
	dexes, err := db.Pokedexes().FindByOwnerID(ctx, user.ID)
	if err != nil {
		return err
	}
	return views.RenderView(c, http.StatusOK, List(dexes))
}

func new(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, New())
}

func create(c echo.Context) error {
	ctx := context.FromEcho(c)
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
	id, err := db.Insert(ctx, dex)
	if err != nil {
		return err
	}

	slog.WithContext(ctx).Infof("created dex %s [%s]", dex.ID, dex.Name)
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/dex/%s", id))
}

func edit(c echo.Context) error {
	c.Set("rendersPokemon", false)
	ctx := context.FromEcho(c)
	id := c.Param("dex")
	db := c.(database.DBContext).DB().Pokedexes()
	dex, err := db.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return views.RenderView(c, http.StatusOK, EditDex(dex))
}

func update(c echo.Context) error {
	ctx := context.FromEcho(c)
	var f struct {
		ID            string             `param:"dex"`
		Name          string             `form:"name"`
		Shiny         string             `form:"shiny"`
		Forms         types.FormLocation `form:"forms"`
		GenderForms   types.FormLocation `form:"gender"`
		RegionalForms types.FormLocation `form:"regional"`
		GMaxForms     types.FormLocation `form:"gmax"`
	}
	if err := c.Bind(&f); err != nil {
		return err
	}

	db := c.(database.DBContext).DB().Pokedexes()
	dex, err := db.FindByID(ctx, f.ID)
	if err != nil {
		return err
	}
	dex.Name = f.Name
	dex.Config = types.PokedexConfig{
		Shiny:         form.ParseBool(f.Shiny),
		Forms:         f.Forms,
		GenderForms:   f.GenderForms,
		RegionalForms: f.RegionalForms,
		GMaxForms:     f.GMaxForms,
	}

	if err = db.Update(ctx, dex); err != nil {
		return err
	}

	slog.WithContext(ctx).Infof("updated dex %s [%s]", dex.ID, dex.Name)
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/dex/%s", f.ID))
}

func delete(c echo.Context) error {
	ctx := context.FromEcho(c)
	id := c.Param("dex")
	db := c.(database.DBContext).DB().Pokedexes()
	if err := db.Delete(ctx, id); err != nil {
		return err
	}

	slog.WithContext(ctx).Infof("deleted dex %s", id)
	return c.Redirect(http.StatusSeeOther, "/dex/")
}
