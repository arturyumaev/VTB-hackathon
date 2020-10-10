package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"payment/domain"
)

type AuthClient struct {
	client   *http.Client
	urlToken string
	secret   string
	logger   *zerolog.Logger
}

func NewAuthClient(urlToken string, secret string, logger *zerolog.Logger) (*AuthClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return &AuthClient{
		client:   client,
		urlToken: urlToken,
		secret:   secret,
		logger:   logger,
	}, nil
}

func (a *AuthClient) CheckAndInvalidateToken(accessToken string, scope []string) (bool, error) {
	tokenReq := domain.TokenRequest{
		AccessToken: accessToken,
		Scope:        scope,
		Secret:       a.secret,
	}

	body, err := json.Marshal(tokenReq)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", a.urlToken, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var tokenResp domain.TokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		return false, err
	}
	return tokenResp.Allowed, nil
}
