package main

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	config2 "payment/config"
	"payment/server"
	"strings"
)

func main() {
	configPathParam := flag.String("config", ".env", "path to config file")
	flag.Parse()

	configPath := strings.TrimSpace(*configPathParam)
	if err := godotenv.Load(configPath); err != nil {
		panic(err)
	}

	var config config2.Config
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	httpServer, err := server.InitHttpServer(config)
	if err != nil {
		panic(err)
	}
	httpServer.StartWebServer()
}
