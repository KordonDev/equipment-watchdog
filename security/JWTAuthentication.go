package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	User               `mapstructure:",squash"`
	jwt.StandardClaims `mapstructure:",squash"`
}

type JwtService struct {
	secretKey string
	issure    string
	domain    string
}

// auth-jwt
func NewJwtService(origin string, jwtSecret, domain string) *JwtService {
	return &JwtService{
		secretKey: jwtSecret,
		issure:    origin,
		domain:    domain,
	}
}

func (service *JwtService) GenerateToken(user User) string {

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

func (service *JwtService) SetCookie(c *gin.Context, token string) {
	c.SetCookie(AUTHORIZATION_COOKIE_KEY, token, 60*100, "/", service.domain, true, true)
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
	User           `mapstructure:",squash"`
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
