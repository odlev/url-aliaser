// Package myauth is a nice package
package myauth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// получаем jwt секретный ключ из енв файла
		jwtSecret, err := secretKey()
		if err != nil {
			log.Error("error loading env files")
		}
		// Получаем токен из заголовка Authorization
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token not provided"})
			log.Error("token not provided")

			return
		}
		// Удаляем "Bearer " из строки (оставляем только сам токен)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Парсим токен и проверяем его подпись
		token, err := jwt.Parse(tokenString, 
			// Проверяем, что алгоритм подписи — HMAC (например, HS256)
			func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
				}
				// Возвращаем секретный ключ для проверки подписи
				return jwtSecret, nil
		})
		// Если токен невалиден (подпись не совпадает или срок истёк)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			log.Error("invalid token")

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["user_id"].(float64) // JWT числа всегда float64!
			c.Set("user_id", int(userID)) // Сохраняем user_id в контекст Gin
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "reading error"})
			log.Error("reading error")

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := int(claims["user_id"].(float64))

			c.Set("user_id", userID)
		}
		// Передаём управление следующему обработчику
		c.Next()
	}
}
