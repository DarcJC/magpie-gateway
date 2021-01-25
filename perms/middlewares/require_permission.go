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


/*
 getValidUser check the Authorization header
 It returns user object and true if success
 nil, false if failed
 */
func getValidUser(c *gin.Context) (*models.AuthorizationUser, bool) {
    reqToken := c.GetHeader("Authorization")
    tokenArr := strings.SplitN(reqToken, "Basic ", 1)
    if len(tokenArr) != 2 {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "bad authorization header",
        })
        return nil, false
    }
    token := tokenArr[1]

    db := store.GetDB()
    session := &models.UserSessionKey{}

    db.Where("key = ? AND is_valid = ?", token, true).First(&session)
    if session.UserID == uuid.Nil || !session.Check() {
        c.JSON(http.StatusUnauthorized, gin.H{
            "code": http.StatusUnauthorized,
            "msg": "bad authorization token",
        })
        db.Updates(session)
        return nil, false
    }

    user := models.AuthorizationUser{}
    db.Where("id = ?", session.UserID).First(&user)

    return &user, true
}


func RequireLoginDecorator(f gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, ok := getValidUser(c)
        if !ok {
            return
        }
        c.Set("user", *user)

        f(c)
    }
}

func RequirePermissionDecorator(f gin.HandlerFunc, perm *perms.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, ok := getValidUser(c)
        if !ok {
            return
        }
        if !perms.TestUserPermission(user, perm.String()) {
            c.JSON(http.StatusForbidden, gin.H{
                "code": http.StatusForbidden,
                "msg": "permission denied",
            })
            return
        }

        c.Set("user", *user)

        f(c)
    }
}
