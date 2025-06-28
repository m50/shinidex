package cmd

import (
	"os"
	"strings"

	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"github.com/m50/shinidex/pkg/maincmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
	Run:   maincmd.Run,
}

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default: './config.yml' or '/config.yml')")
	flags.String("logging.level", "INFO", "log level (default: 'INFO')")
	flags.Bool("logging.access-logs", false, "enabled access logs (default: false)")
	flags.String("logging.format", "text", "format of logs (default: 'text'; options: 'text', 'json')")
	flags.String("listen-address", ":1343", "the address to listen to (default: ':1343')")
	flags.String("db-url", "sqlite://database.db", "the connection url for the database (default: 'sqlite://file:./database.db')")
	flags.BytesHex("auth.key", make([]byte, 32), "they key to be used for signing sessions (default: Regens every run)")
	cobra.CheckErr(viper.BindPFlags(flags))
}

func initConfig() {
	if _, err := os.Stat("./.env"); err == nil {
		slog.Debug("found .env file, loading...")
		if err = godotenv.Load("./.env"); err != nil {
			slog.Errorf("error loading .env file: %s", err)
		}
	}
	viper.SetEnvPrefix("SHINIDEX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		workingDir, err := os.Getwd()
		if err == nil {
			viper.AddConfigPath(workingDir)
		}
		viper.AddConfigPath("/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yml")
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Can't read config:", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Fatal(err)
		os.Exit(1)
	}
}
