package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const ymlConfig = `---
logs:
    format: json
`

func TestInitConfigEnv(t *testing.T) {
	os.Setenv("SHINIDEX_LOGS_FORMAT", "json")
	initConfig()
	assert.Equal(t, "json", viper.GetString("logs.format"))
}

func TestInitConfigYAML(t *testing.T) {
	os.WriteFile("./config.yml", []byte(ymlConfig), 0644)
	defer os.Remove("./config.yml")
	initConfig()
	assert.Equal(t, "json", viper.GetString("logs.format"))
}