package service

import (
	"auth/client"
	"auth/domain"
	"auth/repository"
	"github.com/nu7hatch/gouuid"
	"github.com/rs/zerolog"
)

type Service struct {
	yandexClient *client.YandexClient
	mongo        *repository.MongoDB
	logger       *zerolog.Logger
}

func NewService(mongo *repository.MongoDB, yandexClient *client.YandexClient, logger *zerolog.Logger) *Service {
	return &Service{
		yandexClient: yandexClient,
		mongo:        mongo,
		logger:       logger,
	}
}

func (s *Service) GetConfig(name string) string {
	data, err := s.mongo.GetConfig(name)
	if err != nil {
		return ""
	}
	return data.UrlContent
}

func (s *Service) GetOAuthClient(name string) (*domain.OAuthClient, error) {
	return s.mongo.GetOAuthClient(name)
}

func (s *Service) CheckAndInvalidateUserSession(login string, session string) (bool, error) {
	return s.mongo.CheckAndInvalidateUserSession(login, session)
}

func (s *Service) CheckUserSession(login string, session string) (bool, error) {
	return s.mongo.CheckUserSession(login, session)
}

func (s *Service) InvalidateAllUserSessions(login string) error {
	return s.mongo.InvalidateAllUserSessions(login)
}

func (s *Service) CheckFingerprint(login string, fingerprint string) bool {
	return s.mongo.CheckFingerprint(login, fingerprint)
}

func (s *Service) CheckIp(login string, ip string) bool {
	return s.mongo.CheckIp(login, ip)
}

func (s *Service) WhitelistFingerprint(login string, fingerprint string) error {
	return s.mongo.WhitelistFingerprint(login, fingerprint)
}

func (s *Service) WhitelistIp(login string, ip string) error {
	return s.mongo.WhitelistIp(login, ip)
}

func (s *Service) GenerateAndSaveUserSession(login string) (string, error) {
	random, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	session := random.String()
	return session, s.mongo.SaveUserSession(login, session)
}

func (s *Service) LoadInfo(code string) (domain.UserData, error) {
	data, err := s.mongo.GetOAuthClient("yandex")
	if err != nil {
		return domain.UserData{}, nil
	}
	token, err := s.yandexClient.GetToken(code, data.ClientId, data.ClientPass)
	if err != nil {
		return domain.UserData{}, nil
	}
	userYandexData, err := s.yandexClient.GetUserData(token)
	userData, err := s.mongo.SaveYandexData(userYandexData)
	return userData, err
}

func (s *Service) GetUser(login string, password string) (*domain.UserData, error) {
	return s.mongo.GetUser(login, password)
}

func (s *Service) SaveUser(user domain.UserRegData) error {
	return s.mongo.SaveUser(user)
}