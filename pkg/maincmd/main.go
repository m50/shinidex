package maincmd

import (
	"context"
	"os"
	"time"

	"github.com/gookit/slog"
	"github.com/m50/shinidex/pkg/database"
	imgdownloader "github.com/m50/shinidex/pkg/img-downloader"
	"github.com/m50/shinidex/pkg/logger"
	"github.com/m50/shinidex/pkg/oidc"
	"github.com/m50/shinidex/pkg/search"
	"github.com/m50/shinidex/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	slog.SetLevelByName(viper.GetString("logs.level"))
	logfmt := viper.GetString("logs.format")
	switch logfmt {
	case "text":
		formatter := slog.NewTextFormatter()
		formatter.EnableColor = true
		slog.SetFormatter(formatter)
	case "json":
		slog.SetFormatter(slog.NewJSONFormatter())
	}
	slog.AddProcessor(logger.ContextProcessor{})

	db, err := database.NewFromLoadedConfig()
	if err != nil {
		slog.Fatalf("failed to open DB: %s", err)
		os.Exit(1)
		return
	}
	defer db.Close()
	if err = db.Migrate(); err != nil {
		slog.Fatalf("failed to migrate: %s", err)
		os.Exit(1)
		return
	}

	srch := search.New(db)
	go srch.Prepare(ctx)

	go func() {
		for {
			if err := oidc.Initialize(ctx); err != nil {
				slog.WithContext(ctx).Error(err)
				time.Sleep(2 * time.Second)
				continue
			}
			slog.Info("successfully connected to OIDC provider")
			return
		}
	}()
	go imgdownloader.DownloadImages(ctx, db)

	e := web.New(db, srch)
	if err := e.Start(viper.GetString("listen-address")); err != nil {
		slog.Fatal(err)
		os.Exit(1)
	}
}
