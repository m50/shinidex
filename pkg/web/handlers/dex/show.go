package dex

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/logger"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
)

func showRouter(g *echo.Group) {
	g.GET("/:dex", show)
	g.PATCH("/:dex/pkmn/:pkmn", catch)
	g.DELETE("/:dex/pkmn/:pkmn", release)
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
	entries, err := db.Pokedexes().Entries().List(dex.ID)
	if err != nil {
		return views.RenderError(c, err)
	}

	for i, pkmn := range pokemon {
		idx := slices.IndexFunc(entries, func(e types.PokedexEntry) bool {
			pkmnID, formID := pkmn.IDParts()
			return e.PokemonID == pkmnID && e.FormID == formID
		})
		if idx >= 0 {
			pokemon[i].Caught = true
		}
	}

	return views.RenderView(c, http.StatusOK, Display(pokemon, dex))
}

func catch(c echo.Context) error {
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(c.Param("dex"))
	if err != nil {
		return err
	}
	pkmn, err := db.Pokemon().FindByFullFormID(c.Param("pkmn"))
	if err != nil {
		return err
	}

	pkmnID, formID := pkmn.IDParts()
	logger.Debug(pkmnID, " ", formID)
	if err = db.Pokedexes().Entries().Catch(dex.ID, pkmnID, formID); err != nil {
		return err
	}
	pkmn.Caught = true
	logger.Infof("%s caught for dex %s", pkmn.ID, dex.ID)

	return views.RenderView(c, http.StatusOK,
		Pokemon(dex, pkmn),
		views.Info("Caught", fmt.Sprintf("You caught a %s!", pkmn.Name)),
	)
}

func release(c echo.Context) error {
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(c.Param("dex"))
	if err != nil {
		return err
	}
	pkmn, err := db.Pokemon().FindByFullFormID(c.Param("pkmn"))
	if err != nil {
		return err
	}

	pkmnID, formID := pkmn.IDParts()
	if err = db.Pokedexes().Entries().Release(dex.ID, pkmnID, formID); err != nil {
		return err
	}
	pkmn.Caught = false
	logger.Infof("%s released for dex %s", pkmn.ID, dex.ID)

	return views.RenderView(c, http.StatusOK,
		Pokemon(dex, pkmn),
		views.Info("Released", fmt.Sprintf("You released a %s!", pkmn.Name)),
	)
}
