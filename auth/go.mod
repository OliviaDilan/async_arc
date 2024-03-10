module github.com/OliviaDilan/async_arc/auth

go 1.21.0

require (
	github.com/OliviaDilan/async_arc/pkg v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/ilyakaznacheev/cleanenv v1.5.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/OliviaDilan/async_arc/pkg => ../pkg
