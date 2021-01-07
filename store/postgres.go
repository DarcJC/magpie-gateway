package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgreSQLConnector struct {
	DSN string
}

func (p *PostgreSQLConnector) Connect() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: p.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}
