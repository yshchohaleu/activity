package auth

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
)

type contextKey string

const UserIDKey contextKey = "userID"

func FirebaseAuthMiddleware(app *firebase.App, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No authorization header", http.StatusUnauthorized)
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		auth, err := app.Auth(r.Context())
		if err != nil {
			http.Error(w, "Failed to get Firebase Auth client", http.StatusInternalServerError)
			return
		}

		token, err := auth.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
} 