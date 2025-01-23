package middleware

import (
	"context"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateToken(tokenString)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userIdKey auth.AuthIdKey = "user_id"

		ctx := context.WithValue(r.Context(), userIdKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			http.Error(w, "Guest only", http.StatusForbidden)
		}
		next.ServeHTTP(w, r)
	})
}