package webapp

import (
	"errors"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func getUserIDFromToken(r *http.Request, db *gorm.DB, jwtSecret string) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return 0, errors.New("invalid authorization header format")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrInvalidKey
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, errors.New("email claim not found in token")
	}

	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}