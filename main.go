package main

import (
    "fmt"
    "log"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    _ "magpie-gateway/store"  // init database
)

func main() {
    r := router.SetupRouter()
    if err := r.Run(fmt.Sprintf(":%d", configuration.GlobalConfiguration.Port)); err != nil {
        log.Fatal(err)
    }
}
