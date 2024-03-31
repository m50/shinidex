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
	return nil
}

func toggleCaught(c echo.Context) error {
	ctx := c.(database.DBContext)
	pkmn, err := ctx.DB().Pokemon().FindByID(c.Param("pokemon"))
	if err != nil {
		return views.RenderError(c, err)
	}

	return views.RenderViews(c, http.StatusOK,
		Pokemon(pkmn, true),
		views.Info("Caught", fmt.Sprintf("You caught a %s!", pkmn.Name)),
	)
}
