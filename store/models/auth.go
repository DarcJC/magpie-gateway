package models

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AuthorizationUser struct {
	gorm.Model
	ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username string `gorm:"unique;uniqueIndex"`
	Email string `gorm:"unique"`
	Password string
	Activated bool `gorm:"default:false"`
}

func init() {

}
