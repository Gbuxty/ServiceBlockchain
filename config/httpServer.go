package config

import (
	"time"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Port                  int           `mapstructure:"port"`
	ReadTimeout           time.Duration `mapstructure:"read_timeout"`
	WriteTimeout          time.Duration `mapstructure:"write_timeout"`
	ServerShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func newHTTPServer() HTTPServer {
	return HTTPServer{
		Port:                  viper.GetInt("http.port"),
		ReadTimeout:           viper.GetDuration("http.read_timeout"),
		WriteTimeout:          viper.GetDuration("http.write_timeout"),
		ServerShutdownTimeout: viper.GetDuration("http.shutdown_timeout"),
	}
}
