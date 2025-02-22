package ginerr

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Build is a utility function to build error response.
// It returns JSON in format {"status": "error", "message": message}.
func Build(message string) ErrorResp {
	return ErrorResp{Status: "error", Message: message}
}

// Handle500 deals with server error.
// It is equivalent to c.Error(err) followed by c.JSON(500, Build(err.Error())) call.
func Handle500(c *gin.Context, err error) {
	_ = c.Error(err)
	c.JSON(http.StatusInternalServerError, Build(err.Error()))
}
