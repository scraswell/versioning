package git

import (
	"fmt"

	"github.com/spf13/viper"
)

type VersionBumpStrings struct {
	Major []string `mapstructure:"major"`
	Minor []string `mapstructure:"minor"`
	Patch []string `mapstructure:"patch"`
}

type Config struct {
	BumpStrings      *VersionBumpStrings `mapstructure:"bumpStrings"`
	VersionSeparator string              `mapstructure:"versionSeparator"`
}

func getConfig(v *viper.Viper) *Config {
	myConfig := &Config{}

	err := v.Unmarshal(myConfig)
	if err != nil {
		panic(fmt.Errorf("unable to decode config.  %w", err))
	}

	return myConfig
}

func viperInit() *viper.Viper {
	var viperConfig = viper.New()

	viperConfig.SetConfigName("config")
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath(".")
	viperConfig.AddConfigPath("./config")

	err := viperConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error while reading config: %w", err))
	}

	return viperConfig
}
