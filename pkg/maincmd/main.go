package maincmd

import (
	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/config"
	"github.com/m50/shinidex/pkg/database"
	imgdownloader "github.com/m50/shinidex/pkg/img-downloader"
	l "github.com/m50/shinidex/pkg/logger"
	"github.com/m50/shinidex/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func logger() *slog.Record {
	logger := &slog.Record{}
	slog.SetLevelByName(viper.GetString("log-level"))
	return logger
}

func Run(cmd *cobra.Command, args []string) {
	logger := logger()
	l.SetDefaultLogger(logger)

	db, err := database.NewFromLoadedConfig()
	if err != nil {
		logger.Fatalf("DB failed to open: %s", err)
		return
	}
	defer db.Close()
	if err = db.Migrate("./migrations"); err != nil {
		logger.Fatalf("Failed to migrate: %s", err)
		return
	}

	go imgdownloader.DownloadImages(db)

	e := web.New(db)
	if err := e.Start(config.Loaded.WebAddress); err != nil {
		logger.Fatal(err)
	}
}
