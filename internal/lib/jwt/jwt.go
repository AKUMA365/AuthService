package jwt

import (
	"time"

	"AuthServis/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

func New(user models.User, appSecret string, duration time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
