package main

import (
	"log"

	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/db"
	"github.com/kenmobility/github-api-hex/internal/domain"
)

func main() {
	// load env variables
	config := config.LoadConfig("")

	// establish database connection
	db := db.NewDatabase(*config)

	if err := db.Db.AutoMigrate(&domain.Repository{}, &domain.Commit{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("database migration successful")
}
