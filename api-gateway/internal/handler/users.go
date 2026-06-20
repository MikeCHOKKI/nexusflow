// ─── Users Handlers ─────────────────────────────────────────────
package handler

import (
	"net/http"
	"strconv"

	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
	"github.com/gorilla/mux"
)

// UsersHandler exposes user management endpoints (admin only).
type UsersHandler struct {
	clients *client.Clients
}

func NewUsersHandler(clients *client.Clients) *UsersHandler {
	return &UsersHandler{clients: clients}
}

// GET /api/v1/users
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	resp, err := h.clients.User.ListUsers(r.Context(), &pb.ListUsersRequest{
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list users: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, resp)
}

// GET /api/v1/users/{id}
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := h.clients.User.GetUser(r.Context(), &pb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, user)
}

// PATCH /api/v1/users/{id}/toggle-status
func (h *UsersHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	currentUserID := middleware.UserIDFromContext(r.Context())

	// Prevent self-deactivation
	if id == currentUserID {
		writeError(w, http.StatusBadRequest, "cannot toggle your own status")
		return
	}

	// Get current user to see current status
	user, err := h.clients.User.GetUser(r.Context(), &pb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found: "+err.Error())
		return
	}

	// Toggle is_active
	updated, err := h.clients.User.UpdateUser(r.Context(), &pb.UpdateUserRequest{
		Id:       id,
		IsActive: !user.GetIsActive(),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "toggle status: "+err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, updated)
}
