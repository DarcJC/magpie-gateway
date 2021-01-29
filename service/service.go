package service

import uuid "github.com/satori/go.uuid"

type Type int8

const (
	TypeNone Type = iota
	TypeRest
)

type Base struct {
	ID uuid.UUID
	Type Type
    Endpoints     []Endpoint
}
