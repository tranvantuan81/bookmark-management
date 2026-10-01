package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranvantuan81/bookmark-management/internal/service"
)

const passwordLength = 12

// GenPass is the interface for the genpass handler
type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassService service.GenPass
}

// NewGenPass creates a new genpass handler
func NewGenPass(genPassSvc service.GenPass) GenPass {
	return &genPassHandler{
		genPassService: genPassSvc,
	}
}

// GeneratePassword generates a new password
// @Summary      Generate random password
// @Description  Generates a cryptographically secure random password of fixed length.
// @Tags         password
// @Produce      application/json
// @Success      200  {object}  map[string]string "message"
// @Failure      500  {object}  string "Internal Server Error"
// @Router       /genpass [get]
func (s *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := s.genPassService.GeneratePassword(passwordLength)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"password": pass})
}
