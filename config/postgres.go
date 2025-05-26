package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func newPostgres() PostgresConfig {
	return PostgresConfig{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetInt("postgres.port"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		DBName:   viper.GetString("postgres.db"),
		SSLMode:  viper.GetString("postgres.ssl_mode"),
	}
}

func (p PostgresConfig) ToDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DBName,
		p.SSLMode,
	)
}
