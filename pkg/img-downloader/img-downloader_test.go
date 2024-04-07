package imgdownloader

import (
	"os"
	"sync"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDownloadPokemonImages(t *testing.T) {
	wd, _ := os.Getwd()
	os.RemoveAll(wd + "/assets")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	logger := log.New("test")

	pkmn := types.Pokemon{
		ID: "blastoise",
	}

	if !assert.Contains(t, pkmn.GetImageURL(false), "https://") || !assert.Contains(t, pkmn.GetImageURL(true), "https://") {
		return
	}

	go downloadPokemonImages(wg, pkmn, logger)
	wg.Wait()

	// Make sure image downloaded and saved
	assert.NotContains(t, pkmn.GetImageURL(false), "https://")
	assert.NotContains(t, pkmn.GetImageURL(true), "https://")

	os.RemoveAll(wd + "/assets")
}
