package pokemon

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/views"
)

type Context interface {
	DB() *database.Database
}

func Router(e *echo.Echo) {
	group := e.Group("/pokemon")

	group.GET("/", home).Name = "pokemon-list"
	group.GET("/list", list)
	group.GET("/box/:box", box)
	group.PATCH("/:pokemon", toggleCaught)
}

func home(c echo.Context) error {
	return views.RenderView(c, http.StatusOK, Main())
}

func list(c echo.Context) error {
	ctx := c.(Context)
	pkmn, err := ctx.DB().Pokemon().GetAll()
	if err != nil {
		return views.RenderError(c, err)
	}
	return views.RenderView(c, http.StatusOK, List(pkmn))
}

func box(c echo.Context) error {
	ctx := c.(Context)
	pageNum, err := strconv.Atoi(c.Param("box"))
	if err != nil {
		return views.RenderError(c, err)
	}
	pkmn, err := ctx.DB().Pokemon().Get(30, pageNum)
	if err != nil {
		return views.RenderError(c, err)
	}
	return views.RenderView(c, http.StatusOK, Box(pageNum, pkmn, len(pkmn) == 30))
}

func toggleCaught(c echo.Context) error {
	ctx := c.(Context)
	pkmn, err := ctx.DB().Pokemon().FindByID(c.Param("pokemon"))
	if err != nil {
		return views.RenderError(c, err)
	}
	return views.RenderView(c, http.StatusOK, Pokemon(pkmn, true))
}
