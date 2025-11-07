package main

import (
	"log"

	"github.com/spf13/viper"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Instance              string      `mapstructure:"instance"`
	Addr                  string      `mapstructure:"addr"`
	Username              string      `mapstructure:"username"`
	Password              string      `mapstructure:"password"`
	Database              string      `mapstructure:"database"`
	MaxIdleConnections    int         `mapstructure:"maxIdleConnections"`
	MaxOpenConnections    int         `mapstructure:"maxOpenConnections"`
	MaxConnectionLifeTime string      `mapstructure:"maxConnectionLifeTime"`
	Logger                interface{} `mapstructure:"-"`
}

type Config struct {
	MySQL []*MySQLOptions `mapstructure:"mysql"`
}

func main() {
	v := viper.New()
	v.SetConfigFile("../../configs/config.yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read configs file: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal configs: %v", err)
	}

	log.Printf("Loaded %d MySQL instances\n", len(cfg.MySQL))
	for i, mysqlCfg := range cfg.MySQL {
		log.Printf("Instance %d: %+v\n", i, mysqlCfg)
	}
}
