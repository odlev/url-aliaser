// Package handlers is a nice package
package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/odlev/url-aliaser/internal/config"
	"github.com/odlev/url-aliaser/internal/lib/api"
	"github.com/odlev/url-aliaser/internal/lib/random"
	"github.com/odlev/url-aliaser/internal/lib/sl"
	"github.com/odlev/url-aliaser/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required"` //validate для валидатора
	Alias string `json:"alias,omitempty"`
}

type ResponseAlias struct {
	api.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
	IsAliasExists(alias string) (bool, error)
}

func NewSave(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const operation = "handlers.SaveUrl.NewHandler"

		log = log.With(
			slog.String("operation:", operation),
		)

		// userID := c.MustGet("user_id").(int)

		var req Request

		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode request body"})

			return
		}
		log.Info("request body was decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			
			log.Error("invalid request", sl.Err(err))

			c.JSON(http.StatusBadRequest, api.ValidationError(validateErr))

			return
		}
	
	cfg := config.MustLoad()
	aliasLength := cfg.AliasLength

	alias := req.Alias
	

	if alias == "" {
		for {
			alias = random.NewRandomString(aliasLength)
			exists, err := urlSaver.IsAliasExists(alias)
			if err != nil {
				log.Error("failed to check alias", sl.Err(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "internal error"})
		
				return
			}
			if !exists {
				break
			}
			log.Info("generated alias already exists", slog.String("alias", alias))
		}
	} else {
		exists, err := urlSaver.IsAliasExists(alias)
		if err != nil {
			log.Error("failed to check alias:", sl.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})

			return
		}
		if exists {
			log.Info("alias already exists", slog.String("alias", alias))
			c.JSON(http.StatusConflict, gin.H{"error": "alias already exist"})
			
			return
		}

	} 
	id, err := urlSaver.SaveURL(req.URL, alias)
	if errors.Is(err, storage.ErrURLExist) {
		log.Info("url already exist", slog.String("url", req.URL))
		c.JSON(http.StatusConflict, gin.H{"error": "url already exist", "url": req.URL})

		return
	}
	if err != nil {
		log.Error("failed to add url", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add url"})

		return
	}
	
	log.Info("url succesfully added", slog.Int64("id", id))

	c.JSON(http.StatusCreated, ResponseAlias{Response: api.OK(), Alias: alias})
}
}

