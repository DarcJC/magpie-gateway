package service

import (
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/store"
    "magpie-gateway/store/models"
    "net/http"
    "strings"
)

type Service struct {
	Base
	Source        string     // request source
	directorCache func(r *http.Request)
}

/*
 New return a new instance of service
 param base must have a scheme:
 http://xxxx.xxx:1234/magpie-gateway
 <http|https>://<domain|IP>:<port(number)>[/][path]
 */
func New(uuid uuid.UUID, base string) *Service {
	return &Service{
		Base: Base{
			ID:   uuid,
			Type: TypeRest,
            Endpoints: nil,
		},
		Source:    base,
	}
}

func NewFromModel(service *models.Service) *Service {
    return &Service{
        Base:          Base{
            ID:        service.ID,
            Type:      NewType(service.Info.Type),
            Endpoints: service.Endpoints,
            Token: service.Token,
        },
        Source:        service.Info.Source,
    }
}

/*
 AddEndpoint add a new endpoint to service
 */
func (s *Service) AddEndpoint(endpoint models.ServiceEndpoint) error {
    s.eLock.Lock()
    defer s.eLock.Unlock()

    // check path
    for i := range s.Endpoints {
        if strings.Compare(s.Endpoints[i].Path, endpoint.Path) == 0 {
            return NewError("path existed")
        }
    }

    endpoint.ServiceID = s.ID
    endpoint.ID = 0  // reset id to zero value

    // Add to database
    db := store.GetDB()
    if err := db.Create(&endpoint).Error; err != nil {
        return err
    }

    // Add to slice
    s.Endpoints = append(s.Endpoints, endpoint)

    return nil
}

/*
 LoadEndpoints load service's endpoints into memory
 */
func (s *Service) LoadEndpoints() error {
    var es []models.ServiceEndpoint

    db := store.GetDB()

    s.eLock.Lock()
    defer s.eLock.Unlock()

    if err := db.Where("service_id = ?", s.ID).Preload("Permissions").Find(&es).Error; err != nil {
        return err
    }

    for i := range es {
        AddToRoute(s, &es[i])
        s.Endpoints = append(s.Endpoints, es[i])
    }

    return nil
}

func (s *Service) ReloadEndpoints() error {
    s.eLock.Lock()
    s.Endpoints = make([]models.ServiceEndpoint, 10)
    s.eLock.Unlock()

    return s.LoadEndpoints()
}

/*
 AddPermission create a new permission node for this service
 return ModelError if name == "" or key already exist in this service
*/
func (s *Service) AddPermission(name, desc, key string) error {
    if name == "" {
        return NewError("name could not be an empty string")
    }

    db := store.GetDB()

    var serv models.Service
    if err := db.First(&serv, s.ID).Error; err != nil {
        return err
    }

    if err := db.Preload("Permissions").First(&s, s.ID).Error; err != nil {
        return err
    }

    for i := range serv.Permissions {
        if strings.Compare(serv.Permissions[i].Key, key) == 0 {
            return NewError("exist permission key")
        }
    }

    newPerm := models.PermissionNode{
        ServiceID:   s.ID,
        Key:         key,
        Name:        name,
        Description: desc,
    }

    if err := db.Create(&newPerm).Error; err != nil {
        return err
    }

    return nil
}

func (s *Service) Deactivate() error {
    db := store.GetDB()

    var service models.Service

    if err := db.First(&service, s.ID).Error; err != nil {
        return err
    }

    service.Activated = false

    if err := db.Updates(&service).Error; err != nil {
        return err
    }

    return nil
}

func (s *Service) Activate() error {
    db := store.GetDB()

    var service models.Service

    if err := db.First(&service, s.ID).Error; err != nil {
        return err
    }

    service.Activated = true

    if err := db.Updates(&service).Error; err != nil {
        return err
    }

    return nil
}
