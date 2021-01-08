package configuration

import (
    "github.com/codingconcepts/env"
    "log"
)


type Config struct {
    Port int `env:"MAGPIE_PORT" default:"8080"`
    Debug bool `env:"MAGPIE_DEBUG" default:"true"`
    DatabaseConfig
    SecurityConfig
}

var GlobalConfiguration Config

func init() {
    GlobalConfiguration = Config{}
    if err := env.Set(&GlobalConfiguration); err != nil {
        log.Fatal(err)
    }
    if err := env.Set(&GlobalConfiguration.DatabaseConfig); err != nil {
        log.Fatal(err)
    }
    if err := env.Set(&GlobalConfiguration.SecurityConfig); err != nil {
        log.Fatal(err)
    }
}
