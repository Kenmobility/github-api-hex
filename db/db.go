package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/repositories"
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

// SeedRepository seeds a default chromium repo and set it as tracking
func (d *Database) SeedRepository(config *config.Config) error {
	repository := domain.Repository{
		PublicID:  uuid.New().String(),
		Name:      "chromium/chromium",
		URL:       "https://github.com/chromium/chromium",
		StartDate: config.DefaultStartDate,
		EndDate:   config.DefaultEndDate,
	}

	repoRepository := repositories.NewGormRepositoryRepository(d.Db)

	fRepo, err := repoRepository.RepositoryByName(context.Background(), repository.Name)
	if err == nil {
		//default repository already seeded, set as tracking
		_, err = repoRepository.TrackRepository(context.Background(), *fRepo)
		return err
	}

	sRepo, err := repoRepository.SaveRepository(context.Background(), repository)
	if err != nil {
		return err
	}

	_, err = repoRepository.TrackRepository(context.Background(), *sRepo)
	return err
}
