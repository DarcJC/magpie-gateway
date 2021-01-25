package main

import (
    "fmt"
    "golang.org/x/sync/errgroup"
    "log"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    "magpie-gateway/service/server"
    _ "magpie-gateway/store" // init database
    "net/http"
    "time"
)

var (
    g errgroup.Group
)

func main() {
    server1 := &http.Server{
        Addr: fmt.Sprintf("%s", configuration.GlobalConfiguration.Port),
        Handler: router.SetupRouter(),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    server2 := &http.Server{
        Addr: fmt.Sprintf("%s", configuration.GlobalConfiguration.ServicePort),
        Handler: server.ServiceHandler(),
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    g.Go(func() error {
        err := server1.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
        return err
    })

    g.Go(func() error {
        err := server2.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
        return err
    })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
