package configuration

import (
    "github.com/codingconcepts/env"
    "log"
)


type Config struct {
    Port string `env:"MAGPIE_BIND" default:"localhost:8080"`
    ServicePort string `env:"MAGPIE_SERVICE_BIND" default:"localhost:8081"`
    Debug bool `env:"MAGPIE_DEBUG" default:"true"`
    SID string `env:"MAGPIE_SERVICE_ID" default:"1f6263f3-7b83-4ee8-81a9-c64d3fb251f2"`
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
