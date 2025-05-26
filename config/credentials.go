package config

import "github.com/spf13/viper"

type CredentialsAuth struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}


func newCredentialsAuth() CredentialsAuth {
	return CredentialsAuth{
		User:     viper.GetString("auth.user"),
		Password: viper.GetString("auth.password"),
	}
}
