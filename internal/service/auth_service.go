package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"music-service/internal/models"
	"music-service/internal/repository"
	"time"
)

var (
	invalidSingingMethod   = errors.New("invalid signing method")
	invalidTypeTokenClaims = errors.New("token claims are not of type *tokenClaims")
)

type AuthService struct {
	repo       repository.Authorization
	secret     string
	expiration int64
}

func NewAuthService(
	repo repository.Authorization,
	secret string,
	expiration int64,
) *AuthService {
	return &AuthService{
		repo:       repo,
		secret:     secret,
		expiration: expiration,
	}
}

func (a *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return a.repo.GetUserByEmail(email)
}

func (a *AuthService) CreateUser(user *models.User) error {
	return a.repo.CreateUser(user)
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, invalidSingingMethod
		}

		return []byte(a.secret), nil
	})

	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return 0, invalidTypeTokenClaims
	}

	return claims.UserId, nil
}

func (a *AuthService) CreateToken(username, password string) (string, error) {
	user, err := a.repo.GetUserByUsernameAndPassword(username, password)
	if err != nil {
		return "", err
	}

	expiration := time.Second * time.Duration(a.expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(a.secret))
}
