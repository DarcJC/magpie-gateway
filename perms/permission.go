package perms

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"magpie-gateway/store/models"
)

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
		return checkUserPermission(user, permText)
	}
	return r
}
