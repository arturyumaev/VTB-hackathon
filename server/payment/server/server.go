package server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment/client"
	"payment/config"
	"payment/handlers"
	"payment/repository"
	service2 "payment/service"
	"syscall"
	"time"
)

type HttpServer struct {
	handlerFunc *handlers.HandlerFuncs
	server      *http.Server
	logger      *zerolog.Logger
}

func InitHttpServer(config config.Config) (*HttpServer, error) {
	logger := zerolog.New(os.Stdout)

	mongo, err := repository.NewMongoDB(&config.ConfigMongo, &logger)
	if err != nil {
		panic(err)
	}

	authClient, err := client.NewAuthClient(config.ConfigAuth.TokenUrl, config.ConfigAuth.Secret, &logger)
	if err != nil {
		log.Fatalln(err)
	}

	service := service2.NewService(mongo, authClient, &logger)

	handlerFuncs := handlers.NewHandlerFunc(service, config.JwtKey, &logger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.ConfigServer.Host, config.ConfigServer.Port),
		Handler:      NewRouter(handlerFuncs),
		TLSConfig:    nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	server := &HttpServer{
		handlerFunc: handlerFuncs,
		server:      httpServer,
		logger:      &logger,
	}

	return server, nil
}

func (h *HttpServer) StartWebServer() {
	go func() {
		h.logger.Info().Msg("Start web server")
		_ = h.server.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGILL)
	sig := <-c
	close(c)
	h.stopWebServer()
	h.logger.Info().Msgf("Stop web server, signal = %s", sig)
}

func (h *HttpServer) stopWebServer() {
	if err := h.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
