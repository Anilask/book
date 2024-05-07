package auth

import (
	"context"
	jwt "library/pkg/jwt"
	"net/http"
	"strings"
)

type Claims struct {
	Username string `json:"username"`
	UserType string `json:"userType"`
}

func ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, "Bearer ")
		if len(headerParts) != 2 {
			http.Error(w, "Authorization header must be formatted as 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r)
	}
}
