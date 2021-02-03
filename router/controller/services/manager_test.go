package services

import (
    "bytes"
    "github.com/gin-gonic/gin"
    jsoniter "github.com/json-iterator/go"
    uuid "github.com/satori/go.uuid"
    "github.com/stretchr/testify/assert"
    "magpie-gateway/configuration"
    "magpie-gateway/router"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
    "net/http/httptest"
    "testing"
)

var r *gin.Engine

func init() {
    r = router.SetupRouter()
}

func TestCreateService(t *testing.T) {
    configuration.GlobalConfiguration.DBMock = false

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

    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateExistService(t *testing.T) {
    configuration.GlobalConfiguration.DBMock = false

    w := httptest.NewRecorder()
    data := SMEPutData{
        ID:     uuid.NewV4().String(),
        Name:   "testing",
        Desc:   "a\fgegg esgsfda d",
        Source: "localhost:1080",
    }
    body, _ := jsoniter.Marshal(data)

    req, _ := http.NewRequest("PUT", "/service/manage", bytes.NewReader(body))
    r.ServeHTTP(httptest.NewRecorder(), req)

    req, _ = http.NewRequest("PUT", "/service/manage", bytes.NewReader(body))
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeactivateService(t *testing.T) {
    w := httptest.NewRecorder()

    db := store.GetDB()
    s := models.Service{}
    err := db.First(&s).Error

    assert.Nil(t, err)
    assert.NotNil(t, db)

    data := SMEDeleteData{ID: s.ID.String()}
    body, _ := jsoniter.Marshal(&data)
    req, _ := http.NewRequest("DELETE", "/service/manage", bytes.NewReader(body))
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

}

func TestDeactivateNotExistService(t *testing.T) {
    w := httptest.NewRecorder()

    data := SMEDeleteData{ID: uuid.Nil.String()}
    body, _ := jsoniter.Marshal(&data)
    req, _ := http.NewRequest("DELETE", "/service/manage", bytes.NewReader(body))
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}