package service

import (
	"errors"
	"github.com/rs/zerolog"
	"payment/client"
	"payment/domain"
	"payment/repository"
)

type Service struct {
	authClient *client.AuthClient
	mongo        *repository.MongoDB
	logger       *zerolog.Logger
}

func NewService(mongo *repository.MongoDB, authClient *client.AuthClient, logger *zerolog.Logger) *Service {
	return &Service{
		authClient: authClient,
		mongo:        mongo,
		logger:       logger,
	}
}

func (s *Service) GetBalance(token string, login string) (*domain.UserBalance, error) {
	scope := []string{"balance"}
	if !s.checkAndInvalidateToken(token, scope) {
		return nil, errors.New("no rights")
	}
	data, err := s.mongo.GetBalance(login)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) Pay(token string, login string, addresse string, amount int) error {
	scope := []string{"pay"}
	if !s.checkAndInvalidateToken(token, scope) {
		return errors.New("no rights")
	}
	err := s.mongo.Pay(login, addresse, amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) AddMoney(token string, login string, amount int) error {
	scope := []string{"addMoney"}
	if !s.checkAndInvalidateToken(token, scope) {
		return errors.New("no rights")
	}
	err := s.mongo.AddMoney(login, amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) checkAndInvalidateToken(token string, scope []string) bool {
	isValid, err := s.authClient.CheckAndInvalidateToken(token, scope)
	if err != nil {
		return false
	}
	return isValid
}
