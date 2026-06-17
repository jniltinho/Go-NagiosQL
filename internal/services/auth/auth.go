// Package auth provides JWT-based authentication for NagiosQL.
//
// Password policy:
//   - All new passwords are hashed with bcrypt cost=12 ($2a$12$...).
//   - Legacy MD5 hashes from the PHP project are detected by the absence of the
//     "$2" prefix. A legacy hash blocks login and returns requires_password_reset=true.
//   - No third-party password library is used; only stdlib crypto/md5 and
//     golang.org/x/crypto/bcrypt.
package auth

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrInvalidCredentials is returned when username or password is wrong.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrLegacyMD5 is returned when the stored password is an MD5 hash.
// The caller should inform the client that a password reset is required.
var ErrLegacyMD5 = errors.New("password uses legacy MD5 hash; password reset required")

// ErrUserInactive is returned when the account exists but active='0'.
var ErrUserInactive = errors.New("account is disabled")

// Claims contains the JWT payload for NagiosQL tokens.
type Claims struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	DomainID uint   `json:"domain_id"`
	jwt.RegisteredClaims
}

// TokenPair holds the signed access token and the refresh token string.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Service encapsulates all authentication logic.
type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

// New creates an auth Service.
func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// Login authenticates username+password. On success it returns a TokenPair.
// If the stored hash is MD5, it returns ErrLegacyMD5 without issuing a token.
func (s *Service) Login(username, password string) (*TokenPair, error) {
	var user models.User
	if err := s.db.Where("username = ? AND active = '1'", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("db lookup: %w", err)
	}

	if user.Active == "0" {
		return nil, ErrUserInactive
	}

	// Detect and reject legacy MD5 hashes.
	if isLegacyMD5(user.Password) {
		if !verifyMD5(password, user.Password) {
			return nil, ErrInvalidCredentials
		}
		// Correct password but hash must be upgraded — block login.
		return nil, ErrLegacyMD5
	}

	// Modern bcrypt path.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokenPair(&user)
}

// HashPassword returns a bcrypt hash of password at cost=12.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("bcrypt: %w", err)
	}
	return string(h), nil
}

// UpgradePassword replaces a user's legacy MD5 hash with bcrypt.
// This is called after the user has been verified through a password-reset flow.
func (s *Service) UpgradePassword(userID uint, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("password", hash).Error
}

// ValidateAccessToken parses and validates an access token. Returns its Claims.
func (s *Service) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return claims, nil
}

// RefreshTokens validates a refresh token and issues a fresh TokenPair.
func (s *Service) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Where("username = ? AND active = '1'", claims.Username).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(&user)
}

// issueTokenPair creates a new access+refresh token pair for the given user.
func (s *Service) issueTokenPair(user *models.User) (*TokenPair, error) {
	now := time.Now()

	accessExp := now.Add(time.Duration(s.cfg.JWT.AccessTTLMin) * time.Minute)
	refreshExp := now.Add(time.Duration(s.cfg.JWT.RefreshTTLDays) * 24 * time.Hour)

	makeClaims := func(exp time.Time) Claims {
		return Claims{
			Username: user.Username,
			Admin:    user.Admin == "1",
			DomainID: user.DomainID,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   fmt.Sprintf("%d", user.ID),
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		}
	}

	secret := []byte(s.cfg.JWT.Secret)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, makeClaims(accessExp)).SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, makeClaims(refreshExp)).SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// isLegacyMD5 returns true when hash does not start with the bcrypt prefix "$2".
func isLegacyMD5(hash string) bool {
	return !strings.HasPrefix(hash, "$2")
}

// verifyMD5 checks password against a raw hex-encoded MD5 digest.
func verifyMD5(password, hash string) bool {
	sum := md5.Sum([]byte(password))
	computed := fmt.Sprintf("%x", sum)
	return computed == strings.ToLower(hash)
}
