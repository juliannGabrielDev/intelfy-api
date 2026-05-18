package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/juliannGabrielDev/intelfy-api/pkg/token"
)

// Definimos un tipo personalizado para la clave del Context (Mejor práctica en Go)
type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extraer el header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Token", http.StatusUnauthorized)
			return
		}

		// 2. Formato esperado: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid Token Format", http.StatusUnauthorized)
			return
		}

		// 3. Validar el Token
		userID, role, err := token.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized: Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		// 4. Inyectar el UserID y Role en el Contexto de la petición
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		// 5. Pasar la petición al siguiente handler con el nuevo contexto
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok || userRole != role {
				http.Error(w, "Forbidden: You do not have the required role", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
