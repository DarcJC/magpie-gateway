package perms

import (
    "fmt"
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/configuration"
    "magpie-gateway/store/models"
)

var selfID, _ = uuid.FromString(configuration.GlobalConfiguration.SID)

type Permission struct {
    ServiceID uuid.UUID
    Key string
}

func NewSelfPermission(key string) *Permission {
    return &Permission{
        ServiceID: selfID,
        Key: key,
    }
}

func (p *Permission) String() string {
    return fmt.Sprintf("%s::%s", p.ServiceID, p.Key)
}

func checkUserPermission(user *models.AuthorizationUser, permText string) bool {
    if user.ID != uuid.Nil {
        if user.OwnPerm(permText) {
            return true
        }
        for i := range user.Groups {
            if user.Groups[i].HasPerm(permText) {
                return true
            }
        }
    }
    return false
}

func TestUserPermission(user *models.AuthorizationUser, permText string) bool {
    r, err := getCache(fmt.Sprintf("%s::%s", user.ID, permText))
    if err != nil {
        res := checkUserPermission(user, permText)
        _ = setCache(fmt.Sprintf("%s::%s", user.ID, permText), res)
        return res
    }
    return r
}
