package models

import (
	"fmt"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
)

/*
AuthorizationUser User Model
 */
type AuthorizationUser struct {
	gorm.Model
	ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username string `gorm:"unique;uniqueIndex"`
	Email string `gorm:"unique"`
	Password string
	Activated bool `gorm:"default:false"`
	Groups []PermissionGroup `gorm:"many2many:user_groups;"`
	Permissions []PermissionNode `gorm:"many2many:user_permissions;"`
}

/*
PermissionGroup
 */
type PermissionGroup struct {
	gorm.Model
	ID int `gorm:"primary_key"`
	Name string `gorm:"unique"`
	Description string `gorm:"type:VARCHAR(255)"`
	Title string `gorm:"type:VARCHAR(16)"`
	Permissions []PermissionNode `gorm:"many2many:group_permissions;"`
}

/*
PermissionNode
 */
type PermissionNode struct {
	gorm.Model
	ID int `gorm:"primary_key"`
	ServiceID uuid.UUID `gorm:"type:uuid;"`
	Key string `gorm:"type:VARCHAR(255)"`
	Name string `gorm:"type:VARCHAR(64)"`
	Description string `gorm:"type:VARCHAR(255)"`
}

func (p *PermissionNode) GetPermText() string {
	return fmt.Sprintf("%s::%s", p.ServiceID, p.Key)
}

func (p *PermissionNode) ComparePermText(text string) bool {
	pt := p.GetPermText()
	return strings.Compare(text, pt) == 0
}

func (g *PermissionGroup) HasPerm(text string) bool {
	for i := range g.Permissions {
		if g.Permissions[i].ComparePermText(text) {
			return true
		}
	}
	return false
}

/**
 OwnPerm will not check group permission
 it just check user self owned permissions
 */
func (a *AuthorizationUser) OwnPerm(text string) bool {
	for i := range a.Permissions {
		if a.Permissions[i].ComparePermText(text) {
			return true
		}
	}
	return false
}

func init() {

}
