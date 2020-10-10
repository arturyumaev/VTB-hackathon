package domain

import "github.com/dgrijalva/jwt-go"

type Response struct {
	Status string `json:"status"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

type Claims struct {
	Login string `json:"login"`
	Name string `json:"name"`
	Email string `json:"email"`
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

type TokenRequest struct {
	AccessToken string `json:"accessToken"`
	Scope []string `json:"scope"`
	Secret string `json:"secret"`
}

type TokenResponse struct {
	Allowed bool `json:"allowed"`
}

type UserBalance struct {
	Balance int `json:"balance"`
}

type Transaction struct {
	Login string `bson:"login"`
	Addressee string `bson:"addressee"`
	Anount int `bson:"amount"`
}

type PayRequest struct {
	Addressee string `bson:"addressee"`
	Anount int `bson:"amount"`
}

type AddMoneyRequest struct {
	Anount int `json:"amount"`
}