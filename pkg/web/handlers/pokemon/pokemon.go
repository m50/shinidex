package pokemon

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/search"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/form"
)

func Router(e *echo.Echo) {
	group := e.Group("/pokemon")

	group.GET("", list).Name = "pokemon.list"
	group.GET("/box/:box", box)

	group.GET("/api/search", searchAPI).Name = "pokemon.search"
}

func list(c echo.Context) error {
	ctx := context.FromEcho(c)
	c.Set("rendersPokemon", true)
	db := c.(database.DBContext).DB()
	pkmn, err := db.Pokemon().GetAllAsSeparateForms(ctx)
	if err != nil {
		return views.RenderError(c, err)
	}
	shiny := form.ParseBool(c.FormValue("shiny"))
	return views.RenderView(c, http.StatusOK, List(pkmn, shiny))
}

func box(c echo.Context) error {
	ctx := context.FromEcho(c)
	c.Set("rendersPokemon", true)
	db := c.(database.DBContext).DB()
	pageNum, err := strconv.Atoi(c.Param("box"))
	if err != nil {
		return views.RenderError(c, err)
	}
	pkmn, err := db.Pokemon().GetAllAsSeparateForms(ctx)
	pkmn = pkmn.Box(pageNum - 1)
	if err != nil {
		return views.RenderError(c, err)
	}
	shiny := form.ParseBool(c.FormValue("shiny"))
	return views.RenderView(c, http.StatusOK, Box(pageNum, pkmn, shiny))
}

func searchAPI(c echo.Context) error {
	ctx := context.FromEcho(c)
	s := c.QueryParam("s")
	search := c.(search.Context).Search()
	pkmnID := search.Get(ctx, s)
	return c.JSON(http.StatusOK, map[string]string{"found": pkmnID})
}
