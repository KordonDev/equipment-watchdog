package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type AuthCustomClaims struct {
	Name   string `json:"name"`
	IsUser bool   `json:"isUser"`
	jwt.StandardClaims
}

type JwtService struct {
	secretKey string
	issure    string
}

// auth-jwt
func JWTAuthService(origin string, jwtSecret string) *JwtService {
	return &JwtService{
		secretKey: jwtSecret,
		issure:    origin,
	}
}

func (service *JwtService) GenerateToken(email string, isUser bool) string {
	claims := &AuthCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *JwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})
}

func (service *JwtService) GetClaims(token *jwt.Token) *AuthCustomClaims {
	claims := token.Claims.(jwt.MapClaims)

	var jwtData AuthCustomClaims
	mapstructure.Decode(claims, &jwtData)

	return &jwtData
}
