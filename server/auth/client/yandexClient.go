package client

import (
	"auth/domain"
	"auth/repository"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"net/url"
)

type YandexClient struct {
	client   *http.Client
	urlToken string
	urlLogin string

	logger *zerolog.Logger
}

func NewYandexClient(mongo *repository.MongoDB, logger *zerolog.Logger) (*YandexClient, error) {
	config, err := mongo.GetConfig("yandex_token")
	if err != nil {
		return nil, err
	}
	configLogin, err := mongo.GetConfig("yandex_login")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return &YandexClient{
		client:   client,
		urlToken: config.UrlContent,
		urlLogin: configLogin.UrlContent,
		logger:   logger,
	}, nil
}

func (y *YandexClient) GetToken(code, clientId, clientSecret string) (string, error) {
	formData := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	resp, err := http.PostForm(y.urlToken, formData)
	if err != nil {
		return "", err
	}
	if resp.Body == nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var yandexAnswer domain.YandexToken

	err = json.Unmarshal(body, &yandexAnswer)
	if err != nil {
		return "", err
	}

	if yandexAnswer.Error != "" {
		return fmt.Sprintf("%s, %s", yandexAnswer.Error, yandexAnswer.ErrorDescription), err
	}

	y.logger.Info().Msg("Success got a token")
	return yandexAnswer.AccessToken, nil
}

func (y *YandexClient) GetUserData(token string) (*domain.YandexUser, error) {
	req, err := http.NewRequest("GET", y.urlLogin, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var yandexAnswer domain.YandexUser
	err = json.Unmarshal(body, &yandexAnswer)
	if err != nil {
		return nil, err
	}
	//res, _ := json.Marshal(yandexAnswer)
	return &yandexAnswer, nil
}
