package store

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "magpie-gateway/configuration"
    "magpie-gateway/store/models"
    "sync"
)

type Connector interface {
    Connect() *gorm.DB
}

var defaultConnector Connector

func GetDefaultConnector() *Connector {
    if defaultConnector != nil {
        return &defaultConnector
    }
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
        configuration.GlobalConfiguration.DBHost,
        configuration.GlobalConfiguration.DBUser,
        configuration.GlobalConfiguration.DBPassword,
        configuration.GlobalConfiguration.DBName,
        configuration.GlobalConfiguration.DBPort,
        configuration.GlobalConfiguration.DBSSLMode,
        configuration.GlobalConfiguration.DBTimezone,
    )
    defaultConnector = &PostgreSQLConnector{DSN: dsn}
    return &defaultConnector
}

var dbInstance *gorm.DB
var once sync.Once

func GetDB() *gorm.DB {
    if dbInstance == nil {
        dbInstance = (*GetDefaultConnector()).Connect()
        if configuration.GlobalConfiguration.Debug {
            dbInstance.Logger = logger.Default.LogMode(logger.Info)
        } else {
            dbInstance.Logger = logger.Default.LogMode(logger.Warn)
        }
    }
    return dbInstance
}


func setupDatabase() {
    once.Do(func() {
        db := GetDB()

        db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";") // enable uuid extension to use uuid_generate_v4()

        var err error
        migrate := func(model ...interface{}) {
            if err == nil {
                err = db.AutoMigrate(model...)
            }
        }
        migrate(models.AuthorizationUser{}, models.PermissionNode{}, models.PermissionGroup{}, models.UserSessionKey{})
        migrate(models.Service{})

        if err != nil {
            log.Fatal(err)
        }
    })
}

func init() {
    setupDatabase()
}
