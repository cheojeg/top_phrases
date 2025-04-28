package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	AdminSecretCode      string        `mapstructure:"ADMIN_SECRET_CODE"`
	TwAccessToken        string        `mapstructure:"TW_ACCESS_TOKEN"`
	TwAccessSecret       string        `mapstructure:"TW_ACCESS_SECRET"`
	GotwiApiKey          string        `mapstructure:"GOTWI_API_KEY"`
	GotwiApiKeySecret    string        `mapstructure:"GOTWI_API_KEY_SECRET"`
	PostEnabled          bool          `mapstructure:"POST_ENABLED"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
