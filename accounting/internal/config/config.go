package config

import (
	"github.com/OliviaDilan/async_arc/pkg/config"
)

type Server struct {
	Port string `yaml:"port" env:"PORT"`
	Host string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`

	Auth config.Auth `yaml:"auth"`

	AMQP config.AMQP `yaml:"amqp"`
}