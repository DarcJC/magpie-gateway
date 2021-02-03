package test

import (
    "github.com/gin-gonic/gin"
    "magpie-gateway/router"
    "net/http/httptest"
)

func PrepareTesting() (*gin.Engine, *httptest.ResponseRecorder) {
    return router.SetupRouter(), httptest.NewRecorder()
}
