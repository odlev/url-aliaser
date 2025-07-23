// Package handlers is a nice package
package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odlev/url-aliaser/internal/lib/api"
	"github.com/odlev/url-aliaser/internal/lib/sl"
	"github.com/odlev/url-aliaser/internal/storage"
)

type RequestD struct {
	URL string `json:"url,omitempty"`
	Alias string `json:"alias,omitempty"`
}

type URLDeletter interface {
	DeleteURL(opts storage.DeleteOptions) error
}

func NewDelete(log *slog.Logger, urlDeletter URLDeletter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const operation = "handlers.DeleteURL.NewHandler"

		log = log.With(
			slog.String("operation", operation),
		)
		var reqd RequestD

		err := c.ShouldBindJSON(&reqd)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode request body"})

			return
		}
		log.Info("request body was decoded", slog.Any("request", reqd))

		if err := validator.New().Struct(reqd); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			c.JSON(http.StatusBadRequest, api.ValidationError(validateErr))

			return
		}

		err = urlDeletter.DeleteURL(storage.DeleteOptions{URL: reqd.URL, Alias: reqd.Alias})
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("URL not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})

			return
		}
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Error("alias not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "alias not found"}) 

			return
		}
		if err != nil {
			log.Error("unknown error", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		log.Info("url succesfully deleted")
		c.JSON(http.StatusOK, gin.H{"url": "succesfully deleted"})
	}
}

