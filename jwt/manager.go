package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	errUnauthorized = errors.New("unauthorized")
	errInvalidToken = errors.New("invalid token")
	errTokenExpired = errors.New("token expired")
)

type Manager struct {
	secretKey     []byte
	tokenDuration time.Duration
	keyFunc       func(token *jwt.Token) (interface{}, error)
}

func NewManager(secretKey []byte, tokenDuration time.Duration) *Manager {
	return &Manager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	}
}

func (m *Manager) Generate() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiresAt": &jwt.NumericDate{Time: time.Now().Add(m.tokenDuration)},
	})

	return token.SignedString(m.secretKey)
}

func (m *Manager) Verify(strToken string) error {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(
		strToken,
		&claims,
		m.keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithoutClaimsValidation())

	switch {
	case err != nil:
		return errUnauthorized
	case !token.Valid:
		return errInvalidToken
	case !claims.VerifyExpiresAt(time.Now(), false):
		return errTokenExpired
	default:
		return nil
	}
}
