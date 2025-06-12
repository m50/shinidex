package imgdownloader

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
)

func DownloadImages(db *database.Database, logger *log.Logger) {
	pokemon, err := db.Pokemon().GetAllAsSeparateForms()
	if err != nil {
		logger.Error("Failed to fetch pokemon images: ", err)
		return
	}

	logger.Info("Preparing download of pokemon images...")
	wg := &sync.WaitGroup{}
	wg.Add(len(pokemon))
	for idx, pkmn := range pokemon {
		pkmn := pkmn
		if !strings.Contains(pkmn.GetImageURL(false), "https://") {
			wg.Done()
			continue
		}
		go downloadPokemonImages(wg, pkmn, logger)
		// only do 10 per 5 seconds
		if idx%10 == 0 {
			<-time.After(5 * time.Second)
		}
	}

	wg.Wait()

	logger.Info("Done")
}

func downloadPokemonImages(wg *sync.WaitGroup, pkmn types.Pokemon, logger *log.Logger) {
	// First get normal image
	logger.Debugf("Downloading normal image for %s...", pkmn.ID)
	localPath := pkmn.GetLocalImagePath(false)
	foreignPath := pkmn.GetImageURL(false)
	if err := fetchImage(foreignPath, localPath); err != nil {
		logger.Errorf("Unable to fetch normal local image for pokemon [%s]: %s", pkmn.ID, err)
		wg.Done()
	}
	logger.Infof("Downloaded normal image for %s", pkmn.ID)

	// Now get shiny
	logger.Debugf("Downloading shiny image for %s...", pkmn.ID)
	localPath = pkmn.GetLocalImagePath(true)
	foreignPath = pkmn.GetImageURL(true)
	if err := fetchImage(foreignPath, localPath); err != nil {
		logger.Errorf("Unable to fetch shiny local image for pokemon [%s]: %s", pkmn.ID, err)
		wg.Done()
	}
	logger.Infof("Downloaded shiny image for %s", pkmn.ID)

	// Done
	wg.Done()
}

var httpGet = http.Get

func fetchImage(foreignPath, localPath string) error {
	resp, err := httpGet(foreignPath)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(localPath)
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		return err
	}

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	f.Write(body)
	if err := f.Sync(); err != nil {
		return err
	}
	defer f.Close()

	return nil
}
