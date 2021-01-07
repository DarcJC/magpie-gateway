package auth

import (
	"github.com/gin-gonic/gin"
	"magpie-gateway/router"
	"magpie-gateway/router/controller"
	"net/http"
)

type AuthorizationUserEndpoint struct {
	controller.EndpointBase
}

func (a *AuthorizationUserEndpoint) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": http.StatusOK,
	})
}

func (a *AuthorizationUserEndpoint) Put(c *gin.Context) {
}

func init() {
	endpoint := &AuthorizationUserEndpoint{
		EndpointBase: controller.EndpointBase{
			Path: "test",
		},
	}
	endpoint.Register(endpoint, router.Router.Group("/auth"))
}

