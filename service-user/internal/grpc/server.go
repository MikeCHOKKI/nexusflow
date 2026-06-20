package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	nexusflow "github.com/MikeCHOKKI/nexusflow/service-user/gen"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/db"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/model"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/service"
)

// userServer implémente l'interface nexusflow.UserServiceServer.
type userServer struct {
	nexusflow.UnimplementedUserServiceServer
	auth *service.AuthService
	db   *db.Postgres
}

// NewUserServer crée un nouveau serveur gRPC pour le UserService.
func NewUserServer(authSvc *service.AuthService, database *db.Postgres) nexusflow.UserServiceServer {
	return &userServer{
		auth: authSvc,
		db:   database,
	}
}

// ─── Register ─────────────────────────────────────────────────

func (s *userServer) Register(ctx context.Context, req *nexusflow.RegisterRequest) (*nexusflow.AuthResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if len(req.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	token, expiresAt, user, err := s.auth.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		switch {
		case err == db.ErrAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		default:
			log.Printf("Register error: %v", err)
			return nil, status.Error(codes.Internal, "registration failed")
		}
	}

	return &nexusflow.AuthResponse{
		Token:     token,
		User:      modelToProto(user),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// ─── Login ────────────────────────────────────────────────────

func (s *userServer) Login(ctx context.Context, req *nexusflow.LoginRequest) (*nexusflow.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	token, expiresAt, user, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case err == db.ErrNotFound:
			return nil, status.Error(codes.NotFound, "invalid email or password")
		default:
			log.Printf("Login error: %v", err)
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
	}

	return &nexusflow.AuthResponse{
		Token:     token,
		User:      modelToProto(user),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// ─── GetUser ──────────────────────────────────────────────────

func (s *userServer) GetUser(ctx context.Context, req *nexusflow.GetUserRequest) (*nexusflow.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := s.db.GetUser(ctx, req.Id)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("GetUser error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return modelToProto(user), nil
}

// ─── UpdateUser ───────────────────────────────────────────────

func (s *userServer) UpdateUser(ctx context.Context, req *nexusflow.UpdateUserRequest) (*nexusflow.User, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	fields := make(map[string]any)
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Email != "" {
		fields["email"] = req.Email
	}
	// is_active est un booléen — on l'ajoute toujours (UpdateUserRequest
	// ne permet pas de le laisser non-défini, mais on évite l'écrasement
	// quand il est à false par défaut si l'appelant ne l'a pas modifié).
	// On utilise une astuce: on ne passe is_active que si le caller est admin
	// (à gérer côté gateway) ou on le marque en champ `*bool` dans le proto.
	// Pour simplifier, le proto ayant un bool simple, on le transmet tel quel.
	fields["is_active"] = req.IsActive

	user, err := s.db.UpdateUser(ctx, req.Id, fields)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("UpdateUser error: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return modelToProto(user), nil
}

// ─── ValidateToken ────────────────────────────────────────────

func (s *userServer) ValidateToken(ctx context.Context, req *nexusflow.ValidateTokenRequest) (*nexusflow.User, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	user, err := s.auth.ValidateToken(ctx, req.Token)
	if err != nil {
		log.Printf("ValidateToken error: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	return modelToProto(user), nil
}

// ─── ListUsers ────────────────────────────────────────────────

func (s *userServer) ListUsers(ctx context.Context, req *nexusflow.ListUsersRequest) (*nexusflow.ListUsersResponse, error) {
	users, total, err := s.db.ListUsers(ctx, req.Page, req.Limit, req.Role)
	if err != nil {
		log.Printf("ListUsers error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	protoUsers := make([]*nexusflow.User, len(users))
	for i, u := range users {
		protoUsers[i] = modelToProto(u)
	}

	return &nexusflow.ListUsersResponse{
		Users: protoUsers,
		Pagination: &nexusflow.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		},
	}, nil
}

// ─── Helpers ──────────────────────────────────────────────────

func modelToProto(u *model.User) *nexusflow.User {
	if u == nil {
		return nil
	}
	return &nexusflow.User{
		Id:       u.ID,
		Email:    u.Email,
		Name:     u.Name,
		Role:     u.Role,
		IsActive: u.IsActive,
		Metadata: &nexusflow.Metadata{
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
	}
}


