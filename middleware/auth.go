package middleware

import (
	"errors"
	"net/http"

	"github.com/handymesh/handy_authService/models/session"
	"github.com/handymesh/handy_authService/utils"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var Authorization = r.Header.Get("Authorization")
		if Authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"not auth"`))
			return
		}

		token, err := sessionModel.VerifyToken(Authorization)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"`+err.Error()+`"`))
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"token invalid"`))
		}
		return
	})
}
