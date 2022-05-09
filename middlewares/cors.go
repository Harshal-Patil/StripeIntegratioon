package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		fmt.Println("Access-Control-Allow-Credentials", os.Getenv("CORS_CREDENTIALS"))
		fmt.Println("Access-Control-Allow-Headers", os.Getenv("CORS_HEADERS"))
		fmt.Println("Access-Control-Allow-Methods", os.Getenv("CORS_METHODS"))

		ctx.Header("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		ctx.Header("Access-Control-Allow-Credentials", os.Getenv("CORS_CREDENTIALS"))
		ctx.Header("Access-Control-Allow-Headers", os.Getenv("CORS_HEADERS"))
		ctx.Header("Access-Control-Allow-Methods", os.Getenv("CORS_METHODS"))

		if ctx.Request.Method == "OPTIONS" {
			//TODO: what is our intention of adding this. Please add relevent comments at all such places
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
