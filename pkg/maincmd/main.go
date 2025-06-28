package maincmd

import (
	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/database"
	imgdownloader "github.com/m50/shinidex/pkg/img-downloader"
	"github.com/m50/shinidex/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	slog.SetLevelByName(viper.GetString("logging.level"))
	logfmt := viper.GetString("logging.format")
	switch logfmt {
	case "text":
		slog.SetFormatter(slog.NewTextFormatter())
	case "json":
		slog.SetFormatter(slog.NewJSONFormatter())
	}

	db, err := database.NewFromLoadedConfig()
	if err != nil {
		slog.Fatalf("DB failed to open: %s", err)
		return
	}
	defer db.Close()
	if err = db.Migrate(); err != nil {
		slog.Fatalf("Failed to migrate: %s", err)
		return
	}

	go imgdownloader.DownloadImages(db)

	e := web.New(db)
	if err := e.Start(viper.GetString("listen-address")); err != nil {
		slog.Fatal(err)
	}
}
