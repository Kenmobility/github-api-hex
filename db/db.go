package db

import (
	"github.com/google/uuid"
	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

// NewDatabase creates a connection to db and returns the db instance
func NewDatabase(config config.Config) Database {
	return Database{
		Db: connectPostgresDb(config),
	}
}

type Seeder interface {
	SeedRepository(c *Database, config *config.Config) error
}

// SeedRepository seeds a default chromium repo with tracking as true
func (d *Database) SeedRepository(config *config.Config) error {
	repository := domain.Repository{
		PublicID:   uuid.New().String(),
		Name:       "chromium/chromium",
		URL:        "https://github.com/chromium/chromium",
		IsTracking: true,
		StartDate:  config.DefaultStartDate,
		EndDate:    config.DefaultEndDate,
	}

	err := d.Db.Create(&repository).Error
	if err != nil {
		return err
	}

	return nil
}
