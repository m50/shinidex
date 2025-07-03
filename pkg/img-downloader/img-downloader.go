package imgdownloader

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
)

func DownloadImages(ctx context.Context, db *database.Database) {
	pokemon, err := db.Pokemon().GetAllAsSeparateForms(ctx)
	if err != nil {
		slog.Error("failed to fetch pokemon images: ", err)
		return
	}

	slog.Info("preparing download of pokemon images...")
	wg := &sync.WaitGroup{}
	wg.Add(len(pokemon))
	for idx, pkmn := range pokemon {
		pkmn := pkmn
		if !strings.Contains(pkmn.GetImageURL(false), "https://") {
			wg.Done()
			continue
		}
		go downloadPokemonImages(wg, pkmn)
		// only do 10 per 5 seconds
		if idx%10 == 0 {
			select {
				case <-ctx.Done():
					slog.WithContext(ctx).Warn("context cancel...")
				case <-time.After(5 * time.Second):
			}
		}
	}

	wg.Wait()

	slog.Info("done")
}

func downloadPokemonImages(wg *sync.WaitGroup, pkmn types.Pokemon) {
	// First get normal image
	slog.Debugf("downloading normal image for %s...", pkmn.ID)
	localPath := pkmn.GetLocalImagePath(false)
	foreignPath := pkmn.GetImageURL(false)
	if err := fetchImage(foreignPath, localPath); err != nil {
		slog.Errorf("unable to fetch normal local image for pokemon [%s]: %s", pkmn.ID, err)
		wg.Done()
	}
	slog.Infof("downloaded normal image for %s", pkmn.ID)

	// Now get shiny
	slog.Debugf("downloading shiny image for %s...", pkmn.ID)
	localPath = pkmn.GetLocalImagePath(true)
	foreignPath = pkmn.GetImageURL(true)
	if err := fetchImage(foreignPath, localPath); err != nil {
		slog.Errorf("unable to fetch shiny local image for pokemon [%s]: %s", pkmn.ID, err)
		wg.Done()
	}
	slog.Infof("downloaded shiny image for %s", pkmn.ID)

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
