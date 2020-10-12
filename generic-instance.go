package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type GenericInstance interface {
	GetId() string
}

type GameInstance struct {
	id uuid.UUID
}

func NewGenericInstance() *GameInstance {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Panic(err)
	}

	return &GameInstance{
		id: id,
	}
}

func (instace GameInstance) GetId() string {
	return instace.id.String()
}
