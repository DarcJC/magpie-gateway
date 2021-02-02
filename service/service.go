package service

import (
    uuid "github.com/satori/go.uuid"
    "magpie-gateway/store/models"
    "sync"
)

type Type uint

const (
	TypeUnknown Type = iota
	TypeRest
)

func NewType(u uint) Type {
    switch u {
    case 0:
        return TypeUnknown
    case 1:
        return TypeRest
    default:
        return TypeUnknown
    }
}

type Base struct {
	ID uuid.UUID
	Type Type
    Endpoints     []models.ServiceEndpoint
	eLock sync.Mutex
}
