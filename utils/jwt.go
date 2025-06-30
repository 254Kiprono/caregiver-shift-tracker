package utils

import (
	"caregiver-shift-tracker/config"
	"caregiver-shift-tracker/logger"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var cfg = config.LoadConfig()

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleIDKey contextKey = "role_id"
)

var (
	JWTSecret        = []byte(cfg.JWTSecretKey)
	RefreshJWTSecret = []byte(cfg.JWTRefreshKey)
)

// Claims structure
type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates an access and refresh token
func GenerateJWT(userID int, roleID int) (string, string, error) {
	now := time.Now()
	logger.InfoLogger.Printf("Generating tokens for userID %d, roleID %d at %s", userID, roleID, now.Format(time.RFC3339))

	// Create the access token
	accessTokenClaims := &Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString(JWTSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to sign access token: %v", err)
		return "", "", err
	}

	// Create the refresh token
	refreshTokenClaims := &Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString(RefreshJWTSecret)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to sign refresh token: %v", err)
		return "", "", err
	}

	logger.InfoLogger.Printf("Tokens generated: Access expires at %s, Refresh expires at %s",
		time.Unix(accessTokenClaims.ExpiresAt.Unix(), 0).Format(time.RFC3339),
		time.Unix(refreshTokenClaims.ExpiresAt.Unix(), 0).Format(time.RFC3339))
	return signedAccessToken, signedRefreshToken, nil
}

// ParseToken validates a token and extracts claims
func ParseToken(tokenString string, isRefresh bool) (*Claims, error) {
	secret := JWTSecret
	if isRefresh {
		secret = RefreshJWTSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractJWT parses a JWT token and extracts the user ID and role ID from its claims.
func ExtractJWT(tokenString string, isRefresh bool) (int, int, error) {
	secret := []byte(cfg.JWTSecretKey)
	if isRefresh {
		secret = []byte(cfg.JWTRefreshKey)
	}

	// Parse the token using the same Claims struct
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// Validate the token
	if err != nil || !token.Valid {
		return 0, 0, errors.New("invalid token")
	}

	return claims.UserID, claims.RoleID, nil
}

// AdminOnly ensures the request has a valid JWT and admin access (RoleID = 1)
func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.RespondRaw(ctx, http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, roleID, err := ExtractJWT(token, false)
		if err != nil {
			logger.RespondRaw(ctx, http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			ctx.Abort()
			return
		}

		if roleID != 1 {
			logger.RespondRaw(ctx, http.StatusForbidden, gin.H{"error": "Access denied: Admin role required"})
			ctx.Abort()
			return
		}

		// Pass user info to context
		ctx.Set("user_id", userID)
		ctx.Set("role_id", roleID)
		ctx.Next()
	}
}

func InitJWTConfig(c *config.Config) {
	cfg = c
}
