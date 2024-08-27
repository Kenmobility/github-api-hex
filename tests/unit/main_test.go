package unit

import (
	"fmt"
	"log"
	"testing"

	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func TestMain(m *testing.M) {
	var err error

	config := config.LoadConfig("../../.env")

	conString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)
	fmt.Println("con string: ", conString)
	if helpers.IsLocal() {
		conString += " sslmode=disable"
	}

	db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres database: %v", err)
	}
}
