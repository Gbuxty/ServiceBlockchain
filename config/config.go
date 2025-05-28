package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	pathConfigFile = "dev.yml"
	configType     = "yml"
)

type Config struct {
	HTTP          HTTPServer
	Postgres      PostgresConfig
	BlockchainUrl BlockchainConfig
	Credentials   CredentialsAuth
	Quotes        Quotes
}

func New() (*Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile(pathConfigFile)
	viper.SetConfigType(configType)

	viper.BindEnv("blockchain.url", "BLOCKCHAIN_URL")
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.db", "POSTGRES_DB")
	viper.BindEnv("postgres.ssl_mode", "POSTGRES_SSLMODE")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		HTTP:          newHTTPServer(),
		Postgres:      newPostgres(),
		BlockchainUrl: newBclockchainUrl(),
		Credentials:   newCredentialsAuth(),
		Quotes:        newQuotes(),
	}, nil
}
