package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

// ShortenLink interface for shorten link handler
type ShortenLink interface {
	CreateShortenLink(ctx *gin.Context)
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
	Url string `json:"url"`
	Exp int    `json:"exp"`
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
	input := shortenLinkInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// goi service
	key, err := s.svc.CreateShortenLink(c, input.Url, input.Exp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// tra response
	c.JSON(http.StatusOK, gin.H{"code": key, "message": "Shorten URL generated successfully!"})
}
