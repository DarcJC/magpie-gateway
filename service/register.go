package service

import (
    "errors"
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
)

type Manager struct {
    services map[string]*Service
}

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
    m.services[id] = &Service{
        Base: Base{
            ID:        uid,
            Type:      0,
            Endpoints: nil,
        },
        Source:        source,
    }

    return nil

}
