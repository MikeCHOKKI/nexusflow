// ─── JWT Authentication Middleware ─────────────────────────────
package middleware

import (
	"context"
	"net/http"
	"strings"

	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "user_role"
const UserEmailKey contextKey = "user_email"

// Auth returns an HTTP middleware that validates JWT tokens via UserService.
// It extracts the Bearer token from the Authorization header and calls
// UserService.ValidateToken() over gRPC. On success the user_id and role
// are injected into the request context.
func Auth(clients *client.Clients) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractBearer(r)
			if !ok {
				writeAuthError(w, "missing or malformed authorization header")
				return
			}

			// Forward token to UserService for validation via gRPC.
			ctx := metadata.AppendToOutgoingContext(r.Context(),
				"authorization", "Bearer "+token,
			)
			user, err := clients.User.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
			if err != nil {
				writeAuthError(w, "invalid or expired token: "+err.Error())
				return
			}

			// Inject user context.
			ctx = context.WithValue(r.Context(), UserIDKey, user.GetId())
			ctx = context.WithValue(ctx, UserRoleKey, user.GetRole())
			ctx = context.WithValue(ctx, UserEmailKey, user.GetEmail())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext extracts the authenticated user's ID from the request context.
func UserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// UserRoleFromContext extracts the authenticated user's role from the request context.
func UserRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// UserEmailFromContext extracts the authenticated user's email from the request context.
func UserEmailFromContext(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

func extractBearer(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	if len(auth) < 8 || !strings.EqualFold(auth[:7], "Bearer ") {
		return "", false
	}
	return auth[7:], true
}

func writeAuthError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"success":false,"data":null,"error":"` + msg + `"}`))
}
