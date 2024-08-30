package db

import (
	"fmt"
	"log"

	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgresDb(config config.Config) (*gorm.DB, error) {
	conString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)
	if helpers.IsLocal() {
		conString += " sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to postgres database: %v", err)
		return nil, err
	}
	return db, nil
}
