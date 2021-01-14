package auth

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    "magpie-gateway/router/controller"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
    "time"
)

type LoginUserEndpoint struct {
    controller.EndpointBase
}

type LUEPostData struct {
    Username string `json:"username"`
    Email string `json:"email"`
    Password string `json:"password" binding:"required"`
}

type SessionJWTStruct struct {
    Key string `json:"key"`
    UUID uuid.UUID `json:"uuid"`
    jwt.StandardClaims
}

func (l *LoginUserEndpoint) Post(c *gin.Context) {
    var data LUEPostData
    var user = &models.AuthorizationUser{}

    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "magpie could not bind data from this request",
        })
        return
    }

    if (data.Email == "" && data.Username == "") || data.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "request must include valid username or email and password",
        })
        return
    }

    db := store.GetDB()

    if data.Username != "" {
        db.Where("username = ?", data.Username).First(&user)
    } else {
        db.Where("email = ?", data.Email).First(&user)
    }

    if !checkPassword(user.Password, data.Password) {
        c.JSON(http.StatusForbidden, gin.H{
            "code": http.StatusForbidden,
            "msg": "username or password error",
        })
        return
    }

    sessionKey := models.NewUserSessionKey()
    if err := db.Model(&user).Association("SessionKeys").Append(sessionKey); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "msg": "database error",
        })
        return
    }

    t := time.Now().Unix()
    d, err := time.ParseDuration(configuration.GlobalConfiguration.SessionJWTExpireTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "msg": "server could not parse expire time from env var",
        })
        return
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS384, SessionJWTStruct{
        sessionKey.Key,
        user.ID,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(d).Unix(),
            IssuedAt:  t,
            Issuer:    "magpie",
            NotBefore: t,
            Subject:   "login_session",
        },
    })

    ts, err := token.SignedString(configuration.GlobalConfiguration.SessionJWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code": http.StatusInternalServerError,
            "msg": "could not issue jwt for you",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "code": http.StatusCreated,
        "msg": "session created",
        "data": ts,
    })
}

func init() {
    endpoint := &LoginUserEndpoint{
        EndpointBase: controller.EndpointBase{
            Path: "login",
        },
    }
    endpoint.Register(endpoint, router.Router.Group("/auth"))
}
