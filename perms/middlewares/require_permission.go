package middlewares

import "github.com/gin-gonic/gin"

func RequirePermissionDecorator(f gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {

        f(c)
    }
}
