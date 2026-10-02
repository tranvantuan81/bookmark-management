package handler

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain silences gin's debug output so test failures stay readable.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
