// Package handlers is a nice package
package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odlev/url-aliaser/internal/lib/api"
	"github.com/odlev/url-aliaser/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

type ResponseURL struct {
	api.Response
	URL string
}

type RespErr struct {
	Status string
	Error string
	Err error
}
var statusErr = "Error"

func NewRedirect(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const operation = "handlers.GetURL.NewHandler"

		log = log.With(
			slog.String("operation", operation),
			slog.String("request_id", c.GetString("request_id")),
		)
		

		alias := c.Param("alias")

		if alias == "" {
			log.Info("alias is empty")
			c.JSON(http.StatusBadRequest, RespErr{Status: statusErr, Error: "alias is empty, not found"})

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			c.JSON(http.StatusNotFound, RespErr{Status: statusErr, Error: "url not found"})

			return
		}
		if err != nil {
			log.Info("unknown error, failed to get url")
			c.JSON(http.StatusInternalServerError, RespErr{Status: statusErr, Error: "unknown error", Err: err})

			return
		}
		log.Info("URL succesfully found", slog.String("url", resURL))

		
		c.Redirect(http.StatusFound, resURL)

	}
}
