package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
	"net/http"
	"payment/domain"
	"payment/service"
	"strings"
)

type HandlerFuncs struct {
	service *service.Service
	jwtKey string
	logger  *zerolog.Logger
}

func NewHandlerFunc(service *service.Service, jwtKey string, logger *zerolog.Logger) *HandlerFuncs {
	return &HandlerFuncs{
		service: service,
		jwtKey: jwtKey,
		logger:  logger,
	}
}

func (h *HandlerFuncs) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(token)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	if !contains(claims.Scope, "balance") {
		errorResponse(w, errors.New("no rights"), http.StatusForbidden)
		return
	}
	balance, err := h.service.GetBalance(token, claims.Login)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	userBytes, _ := json.Marshal(balance)
	if userBytes != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(userBytes)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerFuncs) PayHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(token)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	if !contains(claims.Scope, "pay") {
		errorResponse(w, errors.New("no rights"), http.StatusForbidden)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payRequest domain.PayRequest
	err = decoder.Decode(&payRequest)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	err = h.service.Pay(token, claims.Login, payRequest.Addressee, payRequest.Anount)
	if err != nil {
		errorResponse(w, err, http.StatusForbidden)
		return
	}
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) AddMoneyHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(token)
	if err != nil {
		errorResponse(w, err, http.StatusUnauthorized)
		return
	}
	if !contains(claims.Scope, "addMoney") {
		errorResponse(w, errors.New("no rights"), http.StatusForbidden)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var addMoneyRequest domain.AddMoneyRequest
	err = decoder.Decode(&addMoneyRequest)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	err = h.service.AddMoney(token, claims.Login, addMoneyRequest.Anount)
	if err != nil {
		errorResponse(w, err, http.StatusForbidden)
		return
	}
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) decodeJwt(jwtToken string) (*domain.Claims, error) {
	claims := &domain.Claims{}
	tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("token invalid")
	}
	return claims, nil
}

func getToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", errors.New("token not valid or no specified")
	}
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("invalid token parts")
	}
	if strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("invalid token parts")
	}
	return tokenParts[1], nil
}

func writeResponse(w http.ResponseWriter, text string, code int) {
	response := domain.Response{
		Status: text,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(respBytes)
	w.WriteHeader(code)
}

func errorResponse(w http.ResponseWriter, err error, code int) {
	response := domain.ResponseErr{
		Error: err.Error(),
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(respBytes)
	w.WriteHeader(code)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}