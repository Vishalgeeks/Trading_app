package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	AuthContextKey contextKey = "authenticatedUser"
)

type AuthenticatedUser struct {
	UserID string
	Role   string
	JTI    string
}

var jwtSecret string

func SetJWTSecret(secret string) {
	jwtSecret = secret
}

func FromContext(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(AuthContextKey).(*AuthenticatedUser)
	return user, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if jwtSecret == "" {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		authenticatedUser := &AuthenticatedUser{
			UserID: userID,
			Role:   role,
			JTI:    jti,
		}

		ctx := context.WithValue(r.Context(), AuthContextKey, authenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := FromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if user.Role != requiredRole {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
