package config

import (
	"fmt"
)

type Auth struct {
	Host string `yaml:"host" env:"AUTH_HOST"`
	Port string `yaml:"port" env:"AUTH_PORT"`
}

type AMQP struct {
	Host string `yaml:"host" env:"AMQP_HOST"`
	Port string `yaml:"port" env:"AMQP_PORT"`
	User string `yaml:"user" env:"AMQP_USER"`
	Password string `yaml:"password" env:"AMQP_PASSWORD"`
}

func (a AMQP) URI() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", a.User, a.Password, a.Host, a.Port)
}