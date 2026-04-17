package middleware

import (
	"context"
	"master/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func WithJWTAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			http.Error(w, "unautorised", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (any, error) {
			return []byte(config.JWTSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp, err := claims.GetExpirationTime()
			if err != nil || exp == nil || exp.Time.Compare(time.Now()) == -1 {
				// fmt.Println(err,exp.Time)
				http.Error(w, "unautorised", http.StatusUnauthorized)
				return
			}

			user_id, err := claims.GetSubject()
			if user_id == ""{
				http.Error(w, "unautorised", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", user_id)
			h(w, r.WithContext(ctx))
		}else{
				http.Error(w, "unautorised", http.StatusUnauthorized)
				return
		}
	}
}
