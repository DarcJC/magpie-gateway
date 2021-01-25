package service

import uuid "github.com/satori/go.uuid"

type Type int8

const (
	TYPE_NONE Type = iota
	TYPE_REST
)

type Base struct {
	ID uuid.UUID
	Type Type
}
