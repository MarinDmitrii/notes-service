package common

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/MarinDmitrii/notes-service/internal/user/usecase"
)

type contextKey string

const UserContextKey = contextKey("userID")

type AuthMiddleware struct {
	GetUserByEmail *usecase.GetUserByEmailUseCase
}

func NewAuthMiddleware(getUserByEmail *usecase.GetUserByEmailUseCase) *AuthMiddleware {
	return &AuthMiddleware{GetUserByEmail: getUserByEmail}
}

func (a *AuthMiddleware) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		authBase64 := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(authBase64)
		if err != nil {
			http.Error(w, "Invalid Base64 string", http.StatusUnauthorized)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		email, password := credentials[0], credentials[1]

		user, err := a.GetUserByEmail.Execute(r.Context(), email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if password != user.Password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type User struct {
	ID       int
	Email    string
	Password string
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
