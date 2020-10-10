package config

type Config struct {
	ConfigAuth
	ConfigMongo
	ConfigServer
}

type ConfigMongo struct {
	Host     string `env:"PAYMENT_MONGO_HOST" envDefault:"127.0.0.1"`
	Port     string `env:"PAYMENT_MONGO_PORT" envDefault:"27017"`
	Database string `env:"PAYMENT_MONGO_DATABASE" envDefault:"bank"`
}

type ConfigServer struct {
	Host string `env:"PAYMENT_SERVER_HOST" envDefault:"0.0.0.0"`
	Port string `env:"PAYMENT_SERVER_PORT" envDefault:"8081"`
}

type ConfigAuth struct {
	TokenUrl string `env:"PAYMENT_AUTH_TOKEN_URL" envDefault:"https://127.0.0.1:8080/auth/internal/token"`
	Secret string `env:"PAYMENT_AUTH_SECRET" envDefault:"secret"`
	JwtKey string `env:"PAYMENT_AUTH_JWT_KEY" envDefault:"jwtKey"`
}