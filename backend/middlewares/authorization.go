package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RequireAuth(ctx *gin.Context) {
	fmt.Println("require auth")
	ctx.Next()
}
