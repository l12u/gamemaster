package errcode

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Different more detailed error codes.
// Everything below 100 is a `NotFound`.
const (
	GameDoesNotExist   = 000
	PlayerDoesNotExist = 001
	PlayerAlreadyExist = 002
	BoardDoesNotExist  = 003
	InvalidRole        = 101
	InvalidState       = 102
	InvalidType        = 103
)

// D attaches a detailed message and error to the context.
func D(c *gin.Context, status int, msg ...string) {
	c.JSON(getHttpStatusFromCode(status), gin.H{
		"msg":  strings.Join(msg, "\n"),
		"code": status,
	})
}

// S attaches only a simple message
func S(c *gin.Context, httpStatus int, msg ...string) {
	c.JSON(httpStatus, gin.H{
		"msg": strings.Join(msg, "\n"),
	})
}

func getHttpStatusFromCode(status int) (httpStatus int) {
	if status < 100 {
		return http.StatusNotFound
	} else {
		return http.StatusBadRequest
	}
}
