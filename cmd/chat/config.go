package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT"`
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
