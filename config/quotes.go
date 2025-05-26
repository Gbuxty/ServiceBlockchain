package config

import "github.com/spf13/viper"

type Quotes struct {
	Symbols []string `mapstructure:"quotes"`
}

func newQuotes() Quotes {
	return Quotes{
		Symbols: viper.GetStringSlice("quotes"),
	}
}
