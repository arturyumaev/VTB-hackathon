package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"payment/handlers"
)

func NewRouter(h *handlers.HandlerFuncs) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods(http.MethodGet).Path("/payment/balance").HandlerFunc(h.BalanceHandler)
	r.Methods(http.MethodPost).Path("/payment/pay").HandlerFunc(h.PayHandler)
	r.Methods(http.MethodPost).Path("/payment/addMoney").HandlerFunc(h.AddMoneyHandler)
	return r
}
