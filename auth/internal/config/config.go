package config

import "github.com/OliviaDilan/async_arc/pkg/config"

type Server struct {
	Port string `yaml:"port" env:"PORT"`
	Host string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`

	JWT JWT `yaml:"jwt"`

	AMQP config.AMQP `yaml:"amqp"`
}

type JWT struct {
	Secret string `yaml:"secret" env:"JWT_SECRET"`
}