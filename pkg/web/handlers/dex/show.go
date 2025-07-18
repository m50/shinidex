package dex

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
	"github.com/m50/shinidex/pkg/views"
)

func showRouter(g *echo.Group) {
	g.GET("/:dex", show)
	g.PATCH("/:dex/box", boxCatch)
	g.PATCH("/:dex/pkmn/:pkmn", catch)
	g.DELETE("/:dex/pkmn/:pkmn", release)
}

func show(c echo.Context) error {
	ctx := context.FromEcho(c)
	c.Set("rendersPokemon", true)
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(ctx, c.Param("dex"))
	if err != nil {
		return views.RenderError(c, err)
	}
	pokemon, err := db.Pokemon().GetAllAsSeparateForms(ctx)
	if err != nil {
		return views.RenderError(c, err)
	}
	entries, err := db.Pokedexes().Entries().List(ctx, dex.ID)
	if err != nil {
		return views.RenderError(c, err)
	}

	filteredPokemon := slices.Clone(pokemon)
	filteredPokemon = slices.DeleteFunc(filteredPokemon, func(pkmn types.Pokemon) bool {
		if !pkmn.Form {
			return false
		}
		return (!dex.Config.Forms.AfterBaseForm() && pkmn.IsStandardForm()) ||
			(!dex.Config.GMaxForms.AfterBaseForm() && pkmn.IsGMax()) ||
			(!dex.Config.RegionalForms.AfterBaseForm() && pkmn.IsRegional()) ||
			(!dex.Config.GenderForms.AfterBaseForm() && pkmn.IsFemale())
	})
	afterPokemon := types.PokemonList{}
	if dex.Config.Forms.Separate() {
		afterPokemon = append(afterPokemon, pokemon.StandardForms()...)
	}
	if dex.Config.GMaxForms.Separate() {
		afterPokemon = append(afterPokemon, pokemon.GMax()...)
	}
	if dex.Config.GenderForms.Separate() {
		afterPokemon = append(afterPokemon, pokemon.Female()...)
	}
	if dex.Config.RegionalForms.Separate() {
		afterPokemon = append(afterPokemon, pokemon.Regional()...)
	}
	slices.SortStableFunc(afterPokemon, func(pkmnA, pkmnB types.Pokemon) int {
		if pkmnA.NationalDexNumber < pkmnB.NationalDexNumber {
			return -1
		} else if pkmnA.NationalDexNumber > pkmnB.NationalDexNumber {
			return 1
		}
		return 0
	})

	lists := []types.PokemonList{filteredPokemon, afterPokemon}
	for _, list := range lists {
		for i, pkmn := range list {
			idx := slices.IndexFunc(entries, func(e types.PokedexEntry) bool {
				pkmnID, formID := pkmn.IDParts()
				return e.PokemonID == pkmnID && e.FormID == formID
			})
			if idx >= 0 {
				list[i].Caught = true
			}
		}
	}

	return views.RenderView(c, http.StatusOK, Display(lists, dex))
}

func boxCatch(c echo.Context) error {
	ctx := context.FromEcho(c)
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(ctx, c.Param("dex"))
	if err != nil {
		return err
	}
	f := struct {
		Box  int      `json:"box" form:"box"`
		PKMN []string `json:"pkmn" form:"pkmn"`
	}{}
	if err := c.Bind(&f); err != nil {
		return err
	}

	pkmnList := make(types.PokemonList, len(f.PKMN))
	for idx := range f.PKMN {
		pkmn, err := db.Pokemon().FindByFullFormID(ctx, f.PKMN[idx])
		if err != nil {
			return err
		}

		pkmnID, formID := pkmn.IDParts()
		if err = db.Pokedexes().Entries().Catch(ctx, dex.ID, pkmnID, formID); err != nil {
			slog.WithContext(ctx).Error(f.PKMN[idx], ": ", err, "; likely already marked as caught")
		} else {
			slog.WithContext(ctx).Infof("%s: caught for dex %s", pkmn.ID, dex.ID)
		}
		if !pkmn.ShinyLocked {
			pkmn.Caught = true
		}
		pkmnList[idx] = pkmn
	}

	return views.RenderView(c, http.StatusOK,
		Box(dex, f.Box, pkmnList),
		views.Info("Caught", fmt.Sprintf("You caught everything in box %d! Impressive!", f.Box)),
	)
}

func catch(c echo.Context) error {
	ctx := context.FromEcho(c)
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(ctx, c.Param("dex"))
	if err != nil {
		return err
	}
	pkmn, err := db.Pokemon().FindByFullFormID(ctx, c.Param("pkmn"))
	if err != nil {
		return err
	}

	pkmnID, formID := pkmn.IDParts()
	if err = db.Pokedexes().Entries().Catch(ctx, dex.ID, pkmnID, formID); err != nil {
		return err
	}
	if !pkmn.ShinyLocked {
		pkmn.Caught = true
		slog.WithContext(ctx).Infof("%s caught for dex %s", pkmn.ID, dex.ID)
	}

	return views.RenderView(c, http.StatusOK,
		Pokemon(dex, pkmn),
		views.Info("Caught", fmt.Sprintf("You caught a %s!", pkmn.Name)),
	)
}

func release(c echo.Context) error {
	ctx := context.FromEcho(c)
	db := c.(database.DBContext).DB()
	dex, err := db.Pokedexes().FindByID(ctx, c.Param("dex"))
	if err != nil {
		return err
	}
	pkmn, err := db.Pokemon().FindByFullFormID(ctx, c.Param("pkmn"))
	if err != nil {
		return err
	}

	pkmnID, formID := pkmn.IDParts()
	if err = db.Pokedexes().Entries().Release(ctx, dex.ID, pkmnID, formID); err != nil {
		return err
	}
	pkmn.Caught = false
	slog.WithContext(ctx).Infof("%s released for dex %s", pkmn.ID, dex.ID)

	return views.RenderView(c, http.StatusOK,
		Pokemon(dex, pkmn),
		views.Info("Released", fmt.Sprintf("You released a %s!", pkmn.Name)),
	)
}
