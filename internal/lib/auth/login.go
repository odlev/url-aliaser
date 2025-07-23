// Package myauth is a nice package
package myauth

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odlev/url-aliaser/internal/config"
)

type UserRequest struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type passRequst struct {
	user string 
}

func LoginHandler(cfg config.Config, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ureq UserRequest
		c.ShouldBindJSON(&ureq)

		if ureq.User != cfg.User || ureq.Password != cfg.Password {
			log.Error("invalid registration data")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid registration data"})

			return
		}

		var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

		token, err := GenerateToken(seed.Intn(99999999999))
		if err != nil {
			log.Error("error generate token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generate token", "err": err})

			return
		}

		log.Info("token succesfully generated", "token", token)
		c.JSON(http.StatusCreated, gin.H{"token": token})
	}

}
