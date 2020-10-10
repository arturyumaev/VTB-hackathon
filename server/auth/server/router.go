package server

import (
	"auth/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *handlers.HandlerFuncs) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods(http.MethodGet).Path("/auth/profile").HandlerFunc(h.ProfileHandler)
	r.Methods(http.MethodGet).Path("/auth/login/yandex").HandlerFunc(h.LoginYandexHandler)
	r.Methods(http.MethodGet).Path("/auth/return").HandlerFunc(h.ReturnHandler)
	r.Methods(http.MethodPost).Path("/auth/login/self").HandlerFunc(h.LoginSelfHandler)
	r.Methods(http.MethodPost).Path("/auth/logout").HandlerFunc(h.LogoutHandler)
	r.Methods(http.MethodPost).Path("/auth/verify").HandlerFunc(h.VerifyHandler)
	r.Methods(http.MethodPost).Path("/auth/register").HandlerFunc(h.RegisterHandler)
	r.Methods(http.MethodPost).Path("/auth/accessToken").HandlerFunc(h.AccessTokenHandler)
	r.Methods(http.MethodPost).Path("/auth/internal/token").HandlerFunc(h.InternalCheckAndInvalidateTokenHandler)

	r.Methods(http.MethodOptions).Path("/auth/logout").HandlerFunc(h.OptionsHandler)
	r.Methods(http.MethodOptions).Path("/auth/verify").HandlerFunc(h.OptionsHandler)
	r.Methods(http.MethodOptions).Path("/auth/profile").HandlerFunc(h.OptionsHandler)
	r.Methods(http.MethodOptions).Path("/auth/accessToken").HandlerFunc(h.OptionsHandler)
	r.Methods(http.MethodOptions).Path("/auth/login/self").HandlerFunc(h.OptionsHandler)
	r.Methods(http.MethodOptions).Path("/auth/register").HandlerFunc(h.OptionsHandler)
	return r
}
