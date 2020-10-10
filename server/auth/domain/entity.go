package domain

import "github.com/dgrijalva/jwt-go"

type Response struct {
	Status string `json:"status"`
}

type Claims struct {
	Login string `json:"login"`
	Name string `json:"name"`
	Email string `json:"email"`
	Scope []string `json:"scope"`
	Session string `json:"session"`
	jwt.StandardClaims
}

type InternalTokenRequest struct {
	AccessToken string `json:"accessToken"`
	Scope []string `json:"scope"`
	Secret string `json:"secret"`
}


type OAuthConfig struct {
	UrlContent string `json:"urlContent" bson:"urlContent"`
}

type OAuthClient struct {
	Name       string `json:"name" bson:"name"`
	ClientId   string `json:"clientId" bson:"clientId"`
	ClientPass string `json:"clientPass" bson:"clientPass"`
}

type UserData struct {
	Login   string `json:"login" bson:"login"`
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Password   string `json:"-" bson:"passwordHash"`
	Type string `json:"-" bson:"type"`
}

type UserRegData struct {
	Login   string `json:"login" bson:"login"`
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"passwordHash"`
	Type string `json:"-" bson:"type"`
}

type GetTokenRequest struct {
	Scope []string `json:"scope"`
}

type GetTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type CheckTokenRequest struct {
	AccessToken string `json:"accessToken"`
	Scope []string `json:"scope"`
	Secret string `json:"secret"`
}

type TokenResponse struct {
	Allowed bool `json:"allowed"`
}

type LoginRequest struct {
	Login   string `json:"login" bson:"login"`
	Password    string `json:"password" bson:"password"`
}

type Session struct {
	Login   string `json:"login" bson:"login"`
	SessionId    string `json:"sessId" bson:"sessId"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

type YandexToken struct {
	TokenType        string `json:"token_type" bson:"token_type"`
	AccessToken      string `json:"access_token" bson:"access_token"`
	ExpiresIn        int64  `json:"expires_in" bson:"expires_in"`
	RefreshToken     string `json:"refresh_token" bson:"refresh_token"`
	Error            string `json:"error" bson:"error"`
	ErrorDescription string `json:"error_description" bson:"error_description"`
}

type YandexUser struct {
	FirstName    string `json:"first_name" bson:"first_name"`
	LastName     string `json:"last_name" bson:"last_name"`
	Login        string `json:"login" bson:"login"`
	DefaultEmail string `json:"default_email" bson:"default_email"`
	Sex          string `json:"sex" bson:"sex"`
}
