package models

import (
    uuid "github.com/satori/go.uuid"
    "gorm.io/gorm"
    "magpie-gateway/service"
)

type Service struct {
    gorm.Model
    ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
    Name string `gorm:"type:VARCHAR(64)"`
    Description string `gorm:"type:VARCHAR(255);"`
    Permissions []PermissionNode `gorm:"foreignKey:ServiceID"`
    Info ServiceInfo `gorm:"foreignKey:ServiceID"`
    Token uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();"`
    Activated bool `gorm:"default:false"`
}

type ServiceInfo struct {
    gorm.Model
    ServiceID uuid.UUID
    Type service.Type `gorm:"type:INT; default:1;"`
    Source string
}

type ServiceEndpoint struct {
    gorm.Model
    Name string `gorm:"type:VARCHAR(64)"`
    Description string `gorm:"type:VARCHAR(255);"`
    Path string `gorm:"type: VARCHAR(255) NOT NULL;"`
    Permissions []PermissionNode `gorm:"many2many:service_endpoint_permissions;"`
    Activated bool `gorm:"default:false"`
}
