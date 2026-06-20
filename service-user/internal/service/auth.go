package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/MikeCHOKKI/nexusflow/service-user/internal/db"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/model"
)

const (
	tokenExpiry = 24 * time.Hour
	bcryptCost  = 12
)

// AuthService regroupe les opérations d'authentification et de gestion des tokens.
type AuthService struct {
	db    *db.Postgres
	redis *db.RedisClient
	secret []byte
}

// NewAuthService crée un nouveau service d'authentification.
func NewAuthService(database *db.Postgres, redisClient *db.RedisClient, jwtSecret string) *AuthService {
	return &AuthService{
		db:    database,
		redis: redisClient,
		secret: []byte(jwtSecret),
	}
}

// ─── Register ─────────────────────────────────────────────────

// Register crée un compte et retourne un token JWT.
func (s *AuthService) Register(ctx context.Context, email, password, name string) (string, time.Time, *model.User, error) {
	// Vérifier que l'email n'existe pas déjà.
	existing, err := s.db.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return "", time.Time{}, nil, db.ErrAlreadyExists
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return "", time.Time{}, nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.db.CreateUser(ctx, email, hashed, name)
	if err != nil {
		return "", time.Time{}, nil, err
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return "", time.Time{}, nil, fmt.Errorf("generate token: %w", err)
	}

	return token, expiresAt, user, nil
}

// ─── Login ────────────────────────────────────────────────────

// Login authentifie un utilisateur et retourne un token JWT.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, time.Time, *model.User, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, nil, db.ErrNotFound
	}

	if !user.IsActive {
		return "", time.Time{}, nil, fmt.Errorf("account deactivated")
	}

	if err := ComparePassword(user.Password, password); err != nil {
		return "", time.Time{}, nil, fmt.Errorf("invalid credentials")
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return "", time.Time{}, nil, fmt.Errorf("generate token: %w", err)
	}

	return token, expiresAt, user, nil
}

// ─── Token validation ─────────────────────────────────────────

// ValidateToken vérifie et décode un JWT, retourne l'utilisateur correspondant.
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*model.User, error) {
	// Vérifier la blacklist.
	blacklisted, err := s.redis.IsTokenBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("check blacklist: %w", err)
	}
	if blacklisted {
		return nil, fmt.Errorf("token revoked")
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user, err := s.db.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

// ─── Password helpers ─────────────────────────────────────────

// HashPassword retourne le hash bcrypt d'un mot de passe.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword vérifie un mot de passe contre son hash bcrypt.
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ─── JWT internals ────────────────────────────────────────────

type jwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// generateToken crée un JWT signé en HMAC-SHA256 (HS256).
func (s *AuthService) generateToken(user *model.User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(tokenExpiry)

	claims := &jwtClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}
