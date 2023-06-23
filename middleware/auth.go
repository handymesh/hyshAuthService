package middleware

import (
	"errors"
	"net/http"

	sessionModel "github.com/handymesh/hyshAuthService/models/session"
	"github.com/handymesh/hyshAuthService/utils"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var Authorization = r.Header.Get("Authorization")
		if Authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"not auth"`), http.StatusBadRequest)
			return
		}

		token, err := sessionModel.VerifyToken(Authorization)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"token invalid"`), http.StatusBadRequest)
		}
		return
	})
}
