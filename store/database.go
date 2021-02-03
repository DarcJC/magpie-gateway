package store

import (
    "fmt"
    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/postgres"
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

var Mock *sqlmock.Sqlmock

func GetDB() *gorm.DB {
    if configuration.GlobalConfiguration.DBMock && Mock == nil {
        db, mock, err := sqlmock.New()
        Mock = &mock
        if err != nil {
            log.Fatalf("mock failed: %e", err)
        }
        dbInstance, err = gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})
        if err != nil {
            log.Fatalf("gorm mock failed: %e", err)
        }
        setupDatabase()
        return dbInstance
    }
    if dbInstance == nil {
        dbInstance = (*GetDefaultConnector()).Connect()
        if configuration.GlobalConfiguration.Debug {
            dbInstance.Logger = logger.Default.LogMode(logger.Info)
        } else {
            dbInstance.Logger = logger.Default.LogMode(logger.Warn)
        }
        setupDatabase()
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
        migrate(models.ServiceEndpoint{}, models.ServiceInfo{}, models.Service{})

        if err != nil {
            log.Fatal(err)
        }
    })
}

func init() {
}
