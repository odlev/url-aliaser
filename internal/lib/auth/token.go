// Package myauth is a nice package
package myauth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func secretKey() ([]byte, error) {
	if err := godotenv.Load(); err != nil {
		return []byte(""), fmt.Errorf("error loading .env file %w", err)
	}
	jwtSecret := os.Getenv("JWT_KEY")

	return []byte(jwtSecret), nil
}

func GenerateToken(userID int) (string, error) {

	jwtSecret, err := secretKey()
	if err != nil {
		log.Fatal("error loading env files")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(100 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
