package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey         = "your_secret"
	defaultTokenDuration = 24 * time.Hour // 7 days
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, email string, name string) (string, error)
	// ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string) JWTService {

	if secretKey == "" {
		secretKey = jwtSecretKey
	}

	return &jwtService{
		secretKey:     secretKey,
		tokenDuration: defaultTokenDuration,
	}
}
func (js *jwtService) GenerateToken(userId uint, email string, name string) (string, error) {

	// create claims
	claims := JWTClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gotickets",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // create token with claims

	tokenString, err := token.SignedString([]byte(js.secretKey)) // sign token with secret key
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
