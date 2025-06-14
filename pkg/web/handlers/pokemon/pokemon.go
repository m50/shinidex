package pokemon

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/views"
	"github.com/m50/shinidex/pkg/web/form"
)

func Router(e *echo.Echo) {
	group := e.Group("/pokemon")

	group.GET("", list).Name = "pokemon-list"
	group.GET("/box/:box", box)
}

func list(c echo.Context) error {
	c.Set("rendersPokemon", true)
	ctx := c.(database.DBContext)
	pkmn, err := ctx.DB().Pokemon().GetAllAsSeparateForms()
	if err != nil {
		return views.RenderError(c, err)
	}
	shiny := form.ParseBool(c.FormValue("shiny"))
	return views.RenderView(c, http.StatusOK, List(pkmn, shiny))
}

func box(c echo.Context) error {
	c.Set("rendersPokemon", true)
	ctx := c.(database.DBContext)
	pageNum, err := strconv.Atoi(c.Param("box"))
	if err != nil {
		return views.RenderError(c, err)
	}
	// pkmn, err := ctx.DB().Pokemon().Get(30, pageNum)
	pkmn, err := ctx.DB().Pokemon().GetAllAsSeparateForms()
	pkmn = pkmn.Box(pageNum - 1)
	if err != nil {
		return views.RenderError(c, err)
	}
	shiny := form.ParseBool(c.FormValue("shiny"))
	return views.RenderView(c, http.StatusOK, Box(pageNum, pkmn, shiny))
}
