package handlers

import (
	"auth/domain"
	"auth/service"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"time"
)

type HandlerFuncs struct {
	service        *service.Service
	jwtKey         string
	internalSecret         string
	yandexClientId string
	logger         *zerolog.Logger
}

func NewHandlerFunc(service *service.Service, jwtKey string, internalSecret string, logger *zerolog.Logger) (*HandlerFuncs, error) {
	yandexClient, err := service.GetOAuthClient("yandex")
	if err != nil {
		return nil, err
	}
	return &HandlerFuncs{
		service:        service,
		jwtKey:         jwtKey,
		logger:         logger,
		internalSecret: internalSecret,
		yandexClientId: yandexClient.ClientId,
	}, nil
}

func (h *HandlerFuncs) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	var cookie, err = r.Cookie("refreshToken")
	if err != nil || len(cookie.Value) == 0 {
		writeResponse(w, "no auth", http.StatusUnauthorized)
		return
	} else {
		token, err := getToken(r)
		if err != nil || len(token) == 0 {
			writeResponse(w, "no auth", http.StatusUnauthorized)
			return
		}
		claims, err := h.decodeJwt(token, false)
		if err != nil || len(token) == 0 {
			writeResponse(w, "no auth", http.StatusUnauthorized)
			return
		}
		user := domain.UserData{
			Login: claims.Login,
			Name:  claims.Name,
			Email: claims.Email,
		}
		userBytes, _ := json.Marshal(user)
		if userBytes != nil {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(userBytes)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h *HandlerFuncs) LoginYandexHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("client_id", h.yandexClientId)
	loginURL := h.service.GetConfig("yandex_code") + data.Encode()
	redirect := r.URL.Query().Get("redirect_uri")
	cookieReferer := http.Cookie{Name: "redirectUri", Value: redirect, MaxAge: 6000, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookieReferer)

	w.Header().Set("Location", loginURL)
	h.logger.Info().Msg("Go to yandex to take code")
	w.WriteHeader(http.StatusSeeOther)
}

func (h *HandlerFuncs) LoginSelfHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	decoder := json.NewDecoder(r.Body)
	var loginReq domain.LoginRequest
	err := decoder.Decode(&loginReq)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUser(loginReq.Login, loginReq.Password)
	if err != nil || user == nil {
		errorResponse(w, err, http.StatusForbidden)
		return
	}
	refreshToken, err := h.makeToken(*user, []string{"refresh"})
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	cookieSession := http.Cookie{Name: "refreshToken", Value: refreshToken, MaxAge: 6000, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookieSession)
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)

	cookieSession := http.Cookie{Name: "refreshToken", Value: "", MaxAge: 0, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookieSession)
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	decoder := json.NewDecoder(r.Body)
	var registerReq domain.UserRegData
	err := decoder.Decode(&registerReq)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	err = h.service.SaveUser(registerReq)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	decoder := json.NewDecoder(r.Body)
	var verifyReq domain.VerifyRequest
	err := decoder.Decode(&verifyReq)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("refreshToken")
	if err != nil || len(cookie.Value) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}

	token, err := getToken(r)
	if err != nil || len(token) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(token, true)
	if err != nil || len(token) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}
	h.service.WhitelistFingerprint(claims.Login, verifyReq.Analytics.Fingerprint)
	h.service.WhitelistIp(claims.Login, r.RemoteAddr)

	response := domain.VerifyResponse{
		Allowed: true,  // check email code not required in this task
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (h *HandlerFuncs) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	writeResponse(w, "OK", http.StatusOK)
}

func (h *HandlerFuncs) AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	decoder := json.NewDecoder(r.Body)
	var tokenReq domain.GetTokenRequest
	err := decoder.Decode(&tokenReq)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("refreshToken")
	if err != nil || len(cookie.Value) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}

	token, err := getToken(r)
	if err != nil || len(token) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(token, false)
	if err != nil || len(token) == 0 {
		writeResponse(w, "no access", http.StatusUnauthorized)
		return
	}
	user := domain.UserData{
		Login: claims.Login,
		Name:  claims.Name,
		Email: claims.Email,
	}

	fingerprintOk := h.service.CheckFingerprint(claims.Login, tokenReq.Analytics.Fingerprint)
	ipOk := h.service.CheckIp(claims.Login, r.RemoteAddr)

	if !fingerprintOk && !ipOk {
		errorResponse(w, errors.New("verify first"), http.StatusForbidden)
		return
	}
	_, _ = h.decodeJwt(token, true)
	refreshToken, err := h.makeToken(user, []string{"refresh"})
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	accessToken, err := h.makeToken(user, tokenReq.Scope)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	cookieSession := http.Cookie{Name: "refreshToken", Value: refreshToken, MaxAge: 6000, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookieSession)
	response := domain.GetTokenResponse{
		AccessToken: accessToken,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (h *HandlerFuncs) InternalCheckAndInvalidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var tokenReq domain.InternalTokenRequest
	err := decoder.Decode(&tokenReq)
	if err != nil || len(tokenReq.AccessToken) == 0{
		errorResponse(w, err, http.StatusBadRequest)
		return
	}
	if tokenReq.Secret != h.internalSecret {
		writeResponse(w, "no access, invalid secret", http.StatusUnauthorized)
		return
	}
	claims, err := h.decodeJwt(tokenReq.AccessToken, true)
	if err != nil || !contains(claims.Scope, tokenReq.Scope){
		writeResponse(w, "no access, scopes not match", http.StatusUnauthorized)
		return
	}

	response := domain.TokenResponse{
		Allowed: true,
	}
	respBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func (h *HandlerFuncs) ReturnHandler(w http.ResponseWriter, r *http.Request) {
	setCors(&w)
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		h.logger.Info().Msg("Url Param 'code' is missing")
		return
	}
	code := keys[0]

	userData, err := h.service.LoadInfo(code)
	if err != nil {
		writeResponse(w, "error get userData", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.makeToken(userData, []string{"refresh"})
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}
	cookieSession := http.Cookie{Name: "refreshToken", Value: refreshToken, MaxAge: 6000, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookieSession)
	cookie, err := r.Cookie("redirectUri")
	if err != nil {
		writeResponse(w, "OK", http.StatusOK)
	}
	if cookie != nil {
		if len(cookie.Value) > 0 {
			w.Header().Set("Location", cookie.Value)
			w.WriteHeader(http.StatusSeeOther)
		} else {
			writeResponse(w, "OK", http.StatusOK)
		}
	}
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

func (h *HandlerFuncs) makeToken(user domain.UserData, scope []string) (string, error) {
	session, err := h.service.GenerateAndSaveUserSession(user.Login)
	if err != nil {
		return "", err
	}
	claims := domain.Claims{
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
		Scope: scope,
		Session: session,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *HandlerFuncs) decodeJwt(jwtToken string, invalidate bool) (*domain.Claims, error) {
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
	var checked bool
	if invalidate {
		checked, err = h.service.CheckAndInvalidateUserSession(claims.Login, claims.Session)
		if !checked {
		//	_ = h.service.InvalidateAllUserSessions(claims.Login)
		}
	} else {
		checked, err = h.service.CheckUserSession(claims.Login, claims.Session)
	}
	if !checked || err != nil {
		return nil, errors.New("validation not passed, sessions clear")
	}
	return claims, nil
}

func getToken(r *http.Request) (string, error) {
	val, err := r.Cookie("refreshToken")
	if err != nil || val == nil {
		return "", errors.New("no token")
	}
	refreshToken := (*val).Value
	return refreshToken, nil
}

func contains(s []string, e []string) bool {
	var both []string
	for _, a := range s {
		for _, b := range e {
			if a == b {
				both = append(both, a)
			}
		}
	}
	return len(both) == len(e)
}

func setCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}