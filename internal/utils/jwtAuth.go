package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prorok210/AvitoShop/config"
)

func GenerateSecretKey(length int) (string, error) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(key), nil
}

func GenerateAccessToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("Username or email is empty")
	}
	claims := jwt.MapClaims{
		"name":       username,
		"exp":        time.Now().Add(config.JWT_ACCESS_EXPIRATION_TIME).Unix(),
		"token_type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
}

func GenerateRefreshToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("Username or email is empty")
	}
	claims := jwt.MapClaims{
		"name":       username,
		"exp":        time.Now().Add(config.JWT_REFRESH_EXPIRATION_TIME).Unix(),
		"token_type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})

	if err != nil {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		})
		if err != nil {
			return nil, err
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("Token has expired")
			}
		}
		return claims, nil
	}

	return nil, err
}

func RefreshTokens(refreshTokenString string) (string, string, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["token_type"] != "refresh" {
		return "", "", errors.New("Invalid token type")
	}
	username, usernameOk := claims["name"].(string)
	if !usernameOk {
		return "", "", errors.New("Invalid username in token claims")
	}

	newAccessToken, err := GenerateAccessToken(username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(username)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
