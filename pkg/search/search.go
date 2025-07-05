package search

import (
	"context"
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	KeyURL        = "search.meilisearch-url"
	KeyAPIKey     = "search.meilisearch-key"
	KeyAPIKeyFile = "search.meilisearch-key-file"

	IndexPokemon = "shinidex_pokemon"
)

func Flags(flags *pflag.FlagSet) {
	flags.String(KeyURL, "", "the URL of meilisearch instance")
	flags.String(KeyAPIKey, "", "the API key used to access the meilisearch instance")
	flags.String(KeyAPIKeyFile, "", "the path to a file containing the API key used to access the meilisearch instance")
}

func SecretFile() error {
	apiKeySecretFile := viper.GetString(KeyAPIKeyFile)
	if apiKeySecretFile == "" {
		return nil
	}
	apiKey, err := os.ReadFile(apiKeySecretFile)
	if err != nil {
		return err
	}
	viper.Set(KeyAPIKey, apiKey)
	return nil
}

type Context interface {
	Search() *Search
}

type Search struct {
	db    *database.Database
	meili meilisearch.ServiceManager
	ready bool
}

func New(db *database.Database) *Search {
	url := viper.GetString(KeyURL)
	apiKey := viper.GetString(KeyAPIKey)
	var meili meilisearch.ServiceManager
	if url != "" {
		opts := []meilisearch.Option{}
		if apiKey != "" {
			opts = append(opts, meilisearch.WithAPIKey(apiKey))
		}
		meili = meilisearch.New(url, opts...)
	}
	return &Search{
		db:    db,
		meili: meili,
	}
}

func (s *Search) Prepare(ctx context.Context) {
	if s.meili == nil {
		return
	}

	index := s.meili.Index(IndexPokemon)
	pkmn, err := s.db.Pokemon().GetAllAsSeparateForms(ctx)
	if err != nil {
		slog.WithContext(ctx).Errorf("error looking up pokemon, disabling full-text search: %v", err)
		s.meili = nil
		return
	}
	task, err := index.AddDocuments(pkmn)
	if err != nil {
		slog.WithContext(ctx).Errorf("error adding pokemon to full-text search, disabling: %v", err)
		s.meili = nil
		return
	}
	for {
		select {
		case <-ctx.Done():
			slog.Warn("context cancelled")
			return
		case <-time.After(time.Second):
		}
		t, err := index.GetTask(task.TaskUID)
		if err != nil {
			slog.WithContext(ctx).Errorf("error fetching task %v: %v", task.TaskUID, err)
		}
		slog.WithContext(ctx).Debugf("current status of meilisearch document add task %v: %v", task.TaskUID, t.Status)
		if t.Status == meilisearch.TaskStatusSucceeded {
			s.ready = true
			return
		}
		if t.Status == meilisearch.TaskStatusCanceled {
			return
		}
		if t.Status == meilisearch.TaskStatusFailed {
			slog.WithContext(ctx).Errorf("error with task %v: %v", task.TaskUID, t.Error)
			return
		}
	}
}

func (s *Search) Get(ctx context.Context, text string) (pokemonID string) {
	if s.meili == nil || !s.ready || text == "" {
		return text
	}
	resp, err := s.meili.Index(IndexPokemon).SearchWithContext(ctx, text, &meilisearch.SearchRequest{
		AttributesToHighlight: []string{"*"},
		Limit: 1,
	})
	if err != nil {
		slog.WithContext(ctx).Errorf("error searching with meilisearch: %v", err)
		return text
	}

	if len(resp.Hits) < 1 {
		slog.WithContext(ctx).Warn("meilisearch didn't find anything")
		return text
	}

	var pkmn types.Pokemon
	if err := mapstructure.Decode(resp.Hits[0], &pkmn); err != nil {
		slog.WithContext(ctx).Warnf("response didn't provide a pokemon: %v", err)
		return text
	}

	slog.WithContext(ctx).Debugf("found pokemon: %v", pkmn)

	return pkmn.ID
}
