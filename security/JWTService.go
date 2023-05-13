package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kordondev/equipment-watchdog/models"
	"github.com/mitchellh/mapstructure"
)

type StandardClaims struct {
	Audience  string `json:"aud,omitempty" mapstructure:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty" mapstructure:"exp,omitempty"`
	Id        string `json:"jti,omitempty" mapstructure:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty" mapstructure:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty" mapstructure:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty" mapstructure:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty" mapstructure:"sub,omitempty"`
}

type AuthCustomClaims struct {
	models.User        `mapstructure:",squash"`
	jwt.StandardClaims `mapstructure:",squash"`
}

type JwtService struct {
	secretKey string
	issure    string
}

// auth-jwt
func NewJwtService(origin string, jwtSecret string) *JwtService {
	return &JwtService{
		secretKey: jwtSecret,
		issure:    origin,
	}
}

func (service *JwtService) GenerateToken(user models.User) string {

	claims := &AuthCustomClaims{
		user,
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

type Claims struct {
	models.User    `mapstructure:",squash"`
	StandardClaims `mapstructure:",squash"`
}

func (service *JwtService) GetClaims(token *jwt.Token) (Claims, error) {
	var claims Claims
	if err := mapstructure.Decode(token.Claims.(jwt.MapClaims), &claims); err != nil {
		fmt.Printf("Error parsing token: %v", err)
		return Claims{}, err
	}

	return claims, nil
}
