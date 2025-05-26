package config

import "github.com/spf13/viper"

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
	viper.SetConfigFile(pathConfigFile)
	viper.SetConfigType(configType)
	viper.AutomaticEnv()

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
