package config

import "github.com/spf13/viper"

type BlockchainConfig struct {
	URL string `mapstructure:"url"`
}

func newBclockchainUrl() BlockchainConfig {
	return BlockchainConfig{
		URL: viper.GetString("blockchain.url"),
	}
}
