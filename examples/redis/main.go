package main

import (
	"log"

	"github.com/spf13/viper"
)

// RedisOptions defines options for redis.
type RedisOptions struct {
	Instance     string      `mapstructure:"instance"`
	Addrs        []string    `mapstructure:"addrs"`
	Username     string      `mapstructure:"username"`
	Password     string      `mapstructure:"password"`
	DB           int         `mapstructure:"db"`
	PoolSize     int         `mapstructure:"poolSize"`
	MinIdleConns int         `mapstructure:"minIdleConns"`
	Logger       interface{} `mapstructure:"-"`
}

type Config struct {
	Redis []*RedisOptions `mapstructure:"redis"`
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

	log.Printf("Loaded %d Redis instances\n", len(cfg.Redis))
	for i, redisCfg := range cfg.Redis {
		log.Printf("Instance %d: %+v\n", i, redisCfg)
	}
}
