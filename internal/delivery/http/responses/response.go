package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Message string `json:"message,omitempty"`
    Error   string `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
    c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Created(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
    c.JSON(status, Response{Success: false, Error: msg})
}