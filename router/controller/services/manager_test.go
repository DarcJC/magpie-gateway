package services

import (
    "bytes"
    "fmt"
    jsoniter "github.com/json-iterator/go"
    uuid "github.com/satori/go.uuid"
    "github.com/stretchr/testify/assert"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateService(t *testing.T) {
    configuration.GlobalConfiguration.DBMock = false
    r := router.SetupRouter()

    w := httptest.NewRecorder()
    data := SMEPutData{
        ID:     uuid.NewV4().String(),
        Name:   "testing",
        Desc:   "a\fgegg esgsfda d",
        Source: "localhost:1080",
    }
    body, _ := jsoniter.Marshal(data)

    req, _ := http.NewRequest("PUT", "/service/manage", bytes.NewReader(body))
    r.ServeHTTP(w, req)

    fmt.Print(w.Body.String())

    assert.Equal(t, http.StatusCreated, w.Code)
}
