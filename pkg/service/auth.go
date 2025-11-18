package service

import (
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
	repo repository.Autorization
}

func NewAuthService(repo repository.Autorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.PasswordHash = s.generatePasswordHash(user.PasswordHash)
	return s.repo.CreateUser(user)
}
func (s *AuthService) GenerateToken(email, password string) (string, error) { // TEST потім некст
	user, err := s.repo.GetUser(email, s.generatePasswordHash(password))
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
func (s *AuthService) ParseToken(acessToken string) (uint, string, error) {
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
	user, err := s.repo.GetUserByID(claims.UserID)
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
func (s *AuthService) GetUserByID(userID uint) (todo.User, error) {
	return s.repo.GetUserByID(userID)
}
