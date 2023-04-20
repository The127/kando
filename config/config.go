package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Environment string

	Server struct {
		Host string
		Port int

		WriteTimeout time.Duration
		ReadTimeout  time.Duration

		ShutdownWait time.Duration

		MaxReadBytes int64
	}

	Database struct {
		Host     string
		Port     int
		Database string
		User     string
		Password string
		SslMode  string
	}
}

func (c *Config) IsProduction() bool {
	return c.Environment == "Production"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "Staging"
}

func (c *Config) IsDevelopment() bool {
	return !(c.IsStaging() || c.IsProduction())
}

var C Config

func ReadConfig() {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	v.SetEnvPrefix("kando")
	v.AutomaticEnv()

	v.SetConfigFile("config.json")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	setDefaultServerConfig()
	setDefaultDatabaseConfig()

	v.Unmarshal(&C)
}

func setDefaultServerConfig() {
	C.Server.Host = "0.0.0.0"
	C.Server.Port = 8080

	C.Server.ReadTimeout = 15 * time.Second
	C.Server.WriteTimeout = 15 * time.Second

	C.Server.ShutdownWait = 15 * time.Second

	C.Server.MaxReadBytes = 1048576
}

func setDefaultDatabaseConfig() {
	C.Database.SslMode = "disable"
}
