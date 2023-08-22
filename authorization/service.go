package authorization

import (
	"alquran/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Services interface {
	GenerateJWT(uuid, username string) (tokenString, refreshTokenString string, err error)
	ValidateToken(signedToken string) (*jwt.Token, error)
}

type jwtServices struct {
}

var jwtKey = []byte(viper.GetString("SECRET_KEY_JWT"))

type JWTClaim struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewServices() *jwtServices {
	return &jwtServices{}
}

func (s *jwtServices) GenerateJWT(uuid, username string) (tokenString, refreshTokenString string, err error) {
	// Generate access token
	expirationTime := time.Now().Add(24 * 30 * time.Hour)
	uuidHashed, _ := helpers.EncryptString(uuid)
	usernameHashed, _ := helpers.EncryptString(username)

	claims := &JWTClaim{
		Uuid:     uuidHashed,
		Username: usernameHashed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshExpirationTime := time.Now().Add(36 * time.Hour)
	refreshClaims := &JWTClaim{
		Uuid:     uuidHashed,
		Username: usernameHashed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err = refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// Function Validate Token
func (s *jwtServices) ValidateToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, err
}

// Function Generate Refresh Token
func (s *jwtServices) GenerateRefreshToken(uuid, username string) (refreshTokenString string, err error) {
	expirationTime := time.Now().Add(36 * time.Hour) // Refresh token berlaku selama 7 hari
	uuidHashed, _ := helpers.EncryptString(uuid)
	usernameHashed, _ := helpers.EncryptString(username)
	claims := &JWTClaim{
		Uuid:     uuidHashed,
		Username: usernameHashed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err = token.SignedString(jwtKey)
	return
}
