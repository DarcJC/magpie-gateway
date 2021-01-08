package auth

import (
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "golang.org/x/crypto/bcrypt"
    "log"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    "magpie-gateway/router/controller"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
)

func generateHash(pwd string, cost int) string {
    if cost > bcrypt.MaxCost || cost < bcrypt.MinCost {
        cost = bcrypt.DefaultCost
    }
    res, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
    if err != nil {
        log.Println(err)
    }
    return string(res)
}

func checkPassword(hash, pwd string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

type AuthorizationUserEndpoint struct {
    controller.EndpointBase
}

func (a *AuthorizationUserEndpoint) Get(c *gin.Context) {
    c.JSON(200, gin.H{
        "code": http.StatusOK,
    })
}

type AUEPutData struct {
    Username string `json:"username" binding:"required"`
    Email string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func (a *AuthorizationUserEndpoint) Put(c *gin.Context) {
    var user *models.AuthorizationUser
    var data AUEPutData
    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "magpie could not bind data from this request",
        })
        return
    }

    db := store.GetDB()

    user = &models.AuthorizationUser{
        Username:  data.Username,
        Email:     data.Email,
        Password:  generateHash(data.Password, configuration.GlobalConfiguration.EncryptCost),
        Activated: false,
    }
    db.Create(&user)
    if user.ID == uuid.Nil {
        c.JSON(http.StatusConflict, gin.H{
            "code": http.StatusConflict,
            "msg": "user already exist",
        })
        return
    }
    user.Password = ""  // hide password hash
    c.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "msg": "success",
        "data": user,
    })
}

type AUEPatchData struct {
    ID string `json:"id" binding:"required"`
    Activated bool `json:"activated"`
}

func (a *AuthorizationUserEndpoint) Patch(c *gin.Context) {
    var user = &models.AuthorizationUser{}
    var data AUEPatchData
    if err := c.BindJSON(&data); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "magpie could not bind data from this request",
        })
        return
    }

    db := store.GetDB()
    id, err := uuid.FromString(data.ID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code": http.StatusBadRequest,
            "msg": "bad user id",
        })
        return
    }
    db.First(&user, id)
    if user.ID == uuid.Nil {
        c.JSON(http.StatusNotFound, gin.H{
            "code": http.StatusNotFound,
            "msg": "user not found",
        })
        return
    }

    user.Activated = data.Activated
    db.Updates(&user)
    user.Password = ""  // hide password hash
    c.JSON(http.StatusOK, gin.H{
        "code": http.StatusOK,
        "msg": "success",
        "data": user,
    })
}

func init() {
    endpoint := &AuthorizationUserEndpoint{
        EndpointBase: controller.EndpointBase{
            Path: "user",
        },
    }
    endpoint.Register(endpoint, router.Router.Group("/auth"))
}

