package config

type Config struct {
	ConfigMongo
	ConfigServer
}

type ConfigMongo struct {
	Host     string `env:"AUTH_MONGO_HOST" envDefault:"127.0.0.1"`
	Port     string `env:"AUTH_MONGO_PORT" envDefault:"27017"`
	Database string `env:"AUTH_MONGO_DATABASE" envDefault:"mydb"`
}

type ConfigServer struct {
	Host           string `env:"AUTH_SERVER_HOST" envDefault:"0.0.0.0"`
	Port           string `env:"AUTH_SERVER_PORT" envDefault:"8080"`
	JwtKey         string `env:"AUTH_SERVER_JWT_KEY" envDefault:"jwtKey"`
	InternalSecret string `env:"AUTH_INTERNAL_SECRET" envDefault:"secret"`
}
