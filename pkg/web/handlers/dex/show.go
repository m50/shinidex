package dex

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/views"
)

func showRouter(g *echo.Group) {
	g.GET("/:dex", show)
	g.PATCH("/:dex/pkmn/:pkmn", toggleCaught)
}

func show(c echo.Context) error {
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(c.Param("dex"))
	if err != nil {
		return views.RenderError(c, err)
	}
	pokemon, err := db.Pokemon().GetAllAsSeparateForms()
	if err != nil {
		return views.RenderError(c, err)
	}

	return views.RenderView(c, http.StatusOK, Display(pokemon, dex))
}

func toggleCaught(c echo.Context) error {
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(c.Param("dex"))
	if err != nil {
		return views.RenderError(c, err)
	}
	pkmn, err := db.Pokemon().FindByFullFormID(c.Param("pkmn"))
	if err != nil {
		return views.RenderError(c, err)
	}

	return views.RenderView(c, http.StatusOK,
		Pokemon(dex, pkmn, true),
		views.Info("Caught", fmt.Sprintf("You caught a %s!", pkmn.Name)),
	)
}
