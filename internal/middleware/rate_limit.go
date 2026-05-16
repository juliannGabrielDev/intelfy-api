package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// Diccionario para guardar un limitador por cada IP
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

// getVisitor extrae la IP y le asigna un limitador
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Permite 5 peticiones por segundo, con ráfagas de hasta 10
		limiter = rate.NewLimiter(5, 10)
		visitors[ip] = limiter
	}

	return limiter
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener la IP del cliente
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Si el límite se excedió, rechazar la petición
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests - Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
