package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"prac/pkg/repository"
	"prac/todo"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	salt      = "afadmfnaddfaohf9-438fhjahfdkhnadjfb"
	signKey   = "adsbnsvhbdafpskmfpeohnqwpef"
	tockenTTL = 12 * time.Hour
)

type tockenClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user todo.User) (int, error) {
	user.PasswordHash = s.generatePasswordHash(user.PasswordHash)
	return s.repo.CreateUser(ctx, user)
}
func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, email, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	tocken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tockenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tockenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, user.ID,
	})
	return tocken.SignedString([]byte(signKey))
}
func (s *AuthService) ParseToken(ctx context.Context, acessToken string) (uint, string, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tockenClaims{}, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*tockenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tockenClaims")
	}
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get user: %w", err)
	}
	return claims.UserID, user.Role, nil

}
func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (todo.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
