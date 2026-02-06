package usecase

import (
	"errors"
	"technical-test/src/config"
	"technical-test/src/model"
	"technical-test/src/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(name, email, password string) (model.User, error)
	Login(email, password string) (string, model.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
	ErrJWTSecretMissing   = errors.New("jwt secret is not configured")
)

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

func (uc *authUsecase) Register(name, email, password string) (model.User, error) {
	existing, err := uc.userRepo.FindByEmail(email)
	if err == nil && existing.ID != 0 {
		return model.User{}, ErrEmailExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := uc.userRepo.Create(&user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (uc *authUsecase) Login(email, password string) (string, model.User, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", model.User{}, ErrInvalidCredentials
		}
		return "", model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", model.User{}, ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID, user.Email)
	if err != nil {
		return "", model.User{}, err
	}
	return token, user, nil
}

func generateJWT(userID uint, email string) (string, error) {
	if config.JWTSecret == "" {
		return "", ErrJWTSecretMissing
	}

	expMinutes := config.JWTAccessExp
	if expMinutes <= 0 {
		expMinutes = 60
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Duration(expMinutes) * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	if config.JWTSecret == "" {
		return nil, nil, ErrJWTSecretMissing
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	return token, claims, err
}
