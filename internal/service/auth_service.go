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
	authRepo   repository.Authorization
	tokenRepo  repository.Token
	secret     string
	expiration int64
}

func NewAuthService(
	authRepo repository.Authorization,
	secret string,
	expiration int64,
	tokenRepo repository.Token,
) *AuthService {
	return &AuthService{
		authRepo:   authRepo,
		secret:     secret,
		expiration: expiration,
		tokenRepo:  tokenRepo,
	}
}

func (a *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return a.authRepo.GetUserByEmail(email)
}

func (a *AuthService) CreateUser(user *models.User) error {
	return a.authRepo.CreateUser(user)
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
	user, err := a.authRepo.GetUserByUsernameAndPassword(username, password)
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

	signedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	dbToken := models.Token{
		Token:     signedToken,
		ExpiresAt: time.Now().Add(expiration),
		UserID:    user.ID,
	}

	err = a.tokenRepo.SaveToken(dbToken)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthService) InvalidateToken(userID int) error {
	return a.tokenRepo.InvalidateToken(userID)
}

func (a *AuthService) IsTokenValid(token string) (bool, error) {
	return a.tokenRepo.IsTokenValid(token)
}
