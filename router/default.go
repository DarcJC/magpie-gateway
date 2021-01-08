package router

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "magpie-gateway/configuration"
    "magpie-gateway/logger"
    "net/http"
    "time"
)

var Router *gin.Engine

func init() {
    Router = gin.New()
}

func SetupRouter() *gin.Engine {
    r := Router
    r.Use(gin.Recovery())

    if configuration.GlobalConfiguration.Debug {
        gin.SetMode(gin.DebugMode)
        r.Use(gin.Logger())
    } else {
        gin.SetMode(gin.ReleaseMode)
        gin.DisableConsoleColor()

        fileLogger := logger.NewFileLogger(fmt.Sprintf("logs/magpie-%s.log", time.Now().Format("2006-01-02-15-04-05")))
        lc := fileLogger.GetLoggerConfig()
        r.Use(gin.LoggerWithConfig(lc))
        gin.DefaultWriter = lc.Output
    }

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusTeapot, gin.H{
            "code": 200,
            "msg": "pong",
        })
    })

    return r
}
