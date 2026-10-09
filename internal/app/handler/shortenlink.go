package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tranvantuan81/bookmark-management/internal/app/repository"
	"github.com/tranvantuan81/bookmark-management/internal/app/service"
	"github.com/tranvantuan81/bookmark-management/pkg/response"
)

// ShortenLink interface for shorten link handler
type ShortenLink interface {
	CreateShortenLink(ctx *gin.Context)
	Redirect(c *gin.Context)
}

type shortenLink struct {
	svc service.ShortenUrl
}

// NewShortenLink creates a new shorten link handler
func NewShortenLink(svc service.ShortenUrl) ShortenLink {
	return &shortenLink{
		svc: svc,
	}
}

type shortenLinkInput struct {
	Url string `json:"url" binding:"url"`
	Exp int64  `json:"exp" binding:"required,lte=604800"`
}

type shortenLinkRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateShortenLink Generate shorten link
// @Summary      Generate shorten url based on original url that last upto 7 days
// @Description  Generate shorten url based on original url that last upto 7 days
// @Tags         link
// @Accept       application/json
// @Produce      application/json
// @Param        input body shortenLinkInput true "Input required"
// @Success      200  {object} 	map[string]string
// @Router       /v1/links/shorten [post]
func (s *shortenLink) CreateShortenLink(c *gin.Context) {
	// lay input
	input := &shortenLinkInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	// goi service
	key, err := s.svc.CreateShortenLink(c, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create shorten link - Shorten endpoint")
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerError)
		return
	}

	// tra response
	c.JSON(http.StatusOK, gin.H{"code": key, "message": "Shorten URL generated successfully!"})
}

// Redirect Forward the request to the original url
// @Tags link
// @Accept       application/json
// @Produce      application/json
// @Param        code path string true "Shorten link key"
// @Success      302
// @Router       /v1/links/redirect/{code} [get]
func (s *shortenLink) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	url, err := s.svc.GetURLFromCode(c, code)

	if err != nil {
		if errors.Is(err, repository.ErrKeyNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
			return
		}
		log.Error().Err(err).Msg("Failed to get URL from code - Redirect endpoint")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
