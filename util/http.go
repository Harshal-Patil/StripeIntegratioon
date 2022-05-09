package util

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(code int, data interface{}) gin.H {
	return gin.H{
		"response": data,
	}
}
