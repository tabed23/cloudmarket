package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tabed23/cloudmarket-auth/graph/jwt"
)

// AuthMiddleware checks the Authorization header for a JWT and validates it
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Extract token from the header
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr != "" {
				// Validate token and claims
				claims, err := jwt.ValidateJwt(context.Background(), tokenStr)
				if err == nil {
					// Set the claims in the request context
					ctx := context.WithValue(r.Context(), "auth_claims", claims)
					r = r.WithContext(ctx)
				}
			}
		}
		
		next.ServeHTTP(w, r)
	})
}


// CtxValue retrieves JWT claims from the context
func CtxValue(ctx context.Context) *jwt.JwtClaims {
	raw, _ := ctx.Value("auth_claims").(*jwt.JwtClaims)
	return raw
}
