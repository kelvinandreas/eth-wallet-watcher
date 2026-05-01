package service

import (
	"errors"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/config"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/constant"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/model/app"
	"github.com/kelvinandreas/eth-wallet-watcher/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email, password string) error {
	_, err := s.userRepo.FindUserByEmail(email)
	if err == nil {
		return errors.New(constant.ErrEmailAlreadyExists)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(constant.ErrFailedHashPassword)
	}

	newUser := &app.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(newUser)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New(constant.ErrInvalidEmailOrPassword)
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New(constant.ErrInvalidEmailOrPassword)
	}

	claims := &config.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", errors.New(constant.ErrFailedGenerateToken)
	}

	return tokenString, nil
}
