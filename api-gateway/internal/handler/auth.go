// ─── Auth Handlers ─────────────────────────────────────────────
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	clients *client.Clients
}

func NewAuthHandler(clients *client.Clients) *AuthHandler {
	return &AuthHandler{clients: clients}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req pb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := h.clients.User.Register(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "registration failed: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]interface{}{
		"token":      resp.GetToken(),
		"user":       resp.GetUser(),
		"expires_at": resp.GetExpiresAt(),
	})
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	resp, err := h.clients.User.Login(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "login failed: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, map[string]interface{}{
		"token":      resp.GetToken(),
		"user":       resp.GetUser(),
		"expires_at": resp.GetExpiresAt(),
	})
}

// GET /api/v1/auth/me — requires valid JWT.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.clients.User.GetUser(r.Context(), &pb.GetUserRequest{Id: userID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch user: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, user)
}
