package service

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "magpie-gateway/configuration"
    "magpie-gateway/logger"
    "sync"
    "time"
)

type ServerEngine struct {
    Once sync.Once
    Engine *gin.Engine
    Manager *Manager
}

var engine *ServerEngine

func GetServiceEngine() *ServerEngine {
    if engine == nil {
        engine = &ServerEngine{}
        engine.Once.Do(func() {
            engine.Engine = gin.New()
            engine.Engine.Use(gin.Recovery())

            if configuration.GlobalConfiguration.Debug {
                engine.Engine.Use(gin.Logger())
            } else {
                gin.DisableConsoleColor()

                fileLogger := logger.NewFileLogger(fmt.Sprintf("logs/magpie-services-%s.log", time.Now().Format("2006-01-02-15-04-05")))
                lc := fileLogger.GetLoggerConfig()
                engine.Engine.Use(gin.LoggerWithConfig(lc))
                gin.DefaultWriter = lc.Output
            }
            
            engine.Manager = &Manager{}
            if err := engine.Manager.Init(); err != nil {
                log.Fatalf("error while init service manager %e", err)
            }
        })
    }
    return engine
}
