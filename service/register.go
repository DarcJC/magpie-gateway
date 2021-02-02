package service

import (
    "errors"
    "github.com/gin-gonic/gin"
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
    "log"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
    "sync"
    "sync/atomic"
)

type Manager struct {
    services map[string]*Service

    initialized uint32
    reloading uint32
    mu sync.Mutex
}

/*
 GetService return the specific service
 id is the uuid of service
 */
func (m *Manager) GetService(id string) *Service {
    return m.services[id]
}

/*
 CreateService will initialize a new Service and add it into map,
 then create record in database
 If id already exist, return an error and do nothing.
*/
func (m *Manager) CreateService(id, name, desc, source string) error {
    // parse uuid from string BTW valid it
    uid, err := uuid.FromString(id)
    if err != nil {
        // not a uuid, return the error
        return err
    }

    // check services map
    _, ok := m.services[id]
    if ok {
        // found
        return NewError("Service existed in service map")
    }

    db := store.GetDB()
    // check database
    var service models.Service
    err = db.First(&service, id).Error // TODO review required. I don't pretty sure could id be a inject vector.
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        // other errors like connection reset will jump to this branch
        return err
    }

    if err == nil {
        // TODO reload data from database
        return NewError("Service existed in database")
    }

    // create service in database
    // use transaction
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    service = models.Service{
        ID:          uid,
        Name:        name,
        Description: desc,
        Info:        models.ServiceInfo{
            Source: source,
        },
        Activated:   true,
    }
    tx.Create(&service)

    if err := tx.Commit().Error; err != nil {
        return err
    }

    // add service to map
    m.services[id] = New(uid, source)

    return nil

}

func (m *Manager) loadDataFromDB() error {
    var services []models.Service

    db := store.GetDB()
    if res := db.Where("activated = ?", true).Preload("Permissions").Preload("Info").Find(&services); res.Error != nil {
        return res.Error
    }

    for i := range services {
        tmp := NewFromModel(&services[i])
        m.services[tmp.ID.String()] = tmp
        if err := m.services[tmp.ID.String()].LoadEndpoints(); err != nil {
            return err
        }
    }

    return nil
}

/**
 Reload reload data from database
 this action will clear all data in memory
 call when reloading will do nothing
 */
func (m *Manager) Reload() error {
    if atomic.LoadUint32(&m.initialized) == 1 {
        return nil
    }
    atomic.StoreUint32(&m.initialized, 1)

    m.mu.Lock()

    m.services = make(map[string]*Service)
    if err := m.loadDataFromDB(); err != nil {
        return err
    }

    m.mu.Unlock()
    atomic.StoreUint32(&m.initialized, 0)

    return nil
}

func (m *Manager) Init() error {
    if atomic.LoadUint32(&m.initialized) == 1 {
        log.Fatal("duplicate initialization of service handler.")
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    m.services = make(map[string]*Service)
    if err := m.loadDataFromDB(); err != nil {
        return err
    }

    atomic.StoreUint32(&m.initialized, 1)

    return nil
}

func PathExistChecker(handlerFunc gin.HandlerFunc, serviceID string, endpointID uint) gin.HandlerFunc {
    return func(context *gin.Context) {

        e := GetServiceEngine()

        service := e.Manager.GetService(serviceID)

        if service == nil {
            context.JSON(http.StatusNotFound, gin.H{
                "code": http.StatusNotFound,
                "msg": "service unloaded",
            })
            return
        }

        for i := range service.Endpoints {
            if service.Endpoints[i].ID == endpointID {
                handlerFunc(context)
                return
            }
        }

        context.JSON(http.StatusNotFound, gin.H{
            "code": http.StatusNotFound,
            "msg": "endpoint unloaded",
        })
    }
}
