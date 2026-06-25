package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// truoc khi bat dau vao handler (before)
		log.Println("Start func - Check from Middleware")
		ctx.Writer.Write([]byte("Start func - Check from Middleware"))

		ctx.Next()

		// sau khi handler xu ly xong (after)
		log.Println("Start func - End from Middleware")

	}
}
