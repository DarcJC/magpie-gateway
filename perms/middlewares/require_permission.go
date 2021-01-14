package middlewares

import (
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/perms"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
    "strings"
)


func RequireLoginDecorator(f gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        reqToken := c.GetHeader("Authorization")
        tokenArr := strings.SplitN(reqToken, "Basic ", 1)
        if len(tokenArr) != 2 {
            c.JSON(http.StatusBadRequest, gin.H{
                "code": http.StatusBadRequest,
                "msg": "bad authorization header",
            })
            return
        }
        token := tokenArr[1]

        db := store.GetDB()
        session := &models.UserSessionKey{}

        db.Where("key = ?", token).First(&session)
        if session.UserID == uuid.Nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": http.StatusUnauthorized,
                "msg": "bad authorization token",
            })
            return
        }

        user := models.AuthorizationUser{}
        db.Where("id = ?", session.UserID).First(&user)
        c.Set("user", user)

        f(c)
    }
}

type Permission struct {
}

func RequirePermissionDecorator(f gin.HandlerFunc, permText string) gin.HandlerFunc {
    return func(c *gin.Context) {
        reqToken := c.GetHeader("Authorization")
        tokenArr := strings.SplitN(reqToken, "Basic ", 1)
        if len(tokenArr) != 2 {
            c.JSON(http.StatusBadRequest, gin.H{
                "code": http.StatusBadRequest,
                "msg": "bad authorization header",
            })
            return
        }
        token := tokenArr[1]

        db := store.GetDB()
        session := &models.UserSessionKey{}

        db.Where("key = ?", token).First(&session)
        if session.UserID == uuid.Nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": http.StatusUnauthorized,
                "msg": "bad authorization token",
            })
            return
        }

        user := models.AuthorizationUser{}
        db.Where("id = ?", session.UserID).First(&user)
        if !perms.TestUserPermission(&user, permText) {
            c.JSON(http.StatusForbidden, gin.H{
                "code": http.StatusForbidden,
                "msg": "permission denied",
            })
            return
        }

        c.Set("user", user)

        f(c)
    }
}
