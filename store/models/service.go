package models

import (
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
)

type Service struct {
    gorm.Model
    ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
    Name string `gorm:"type:VARCHAR(64)"`
    Path string `gorm:"type:VARCHAR(64)`
    Description string `gorm:"type:VARCHAR(255);"`
    Permissions []PermissionNode `gorm:"foreignKey:ServiceID"`
    Token uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();"`
    Activated bool `gorm:"default:false"`
}
