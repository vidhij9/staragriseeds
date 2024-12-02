package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/pkg/auth"
	"backend/pkg/errors"
)

var userId = "user_id"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.WriteJSONError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.WriteJSONError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		claims, err := auth.VerifyToken(token)
		if err != nil {
			errors.WriteJSONError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userId, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
