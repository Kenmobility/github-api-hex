package domain

import (
	"context"
	"time"

	"github.com/kenmobility/github-api-hex/dtos"
)

type Repository struct {
	ID              uint   `gorm:"primarykey"`
	PublicID        string `gorm:"type:varchar;uniqueIndex"`
	Name            string `gorm:"type:varchar;unique"`
	Description     string `gorm:"type:text"`
	URL             string `gorm:"type:varchar"`
	Language        string `gorm:"type:varchar"`
	ForksCount      int
	StarsCount      int
	OpenIssuesCount int
	WatchersCount   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartDate       time.Time
	EndDate         time.Time
	IsTracking      bool
}

type RepositoryRepository interface {
	SaveRepository(ctx context.Context, repository Repository) (*Repository, error)
	RepositoryByPublicId(ctx context.Context, publicId string) (*Repository, error)
	RepositoryByName(ctx context.Context, name string) (*Repository, error)
	AllRepositories(ctx context.Context) ([]Repository, error)
	TrackedRepository(ctx context.Context) (*Repository, error)
	TrackRepository(ctx context.Context, repository Repository) (*Repository, error)
}

func (r Repository) ToDto() dtos.RepositoryResponse {
	return dtos.RepositoryResponse{
		Id:              r.PublicID,
		Name:            r.Name,
		Description:     r.Description,
		URL:             r.URL,
		Language:        r.Language,
		ForksCount:      r.ForksCount,
		StarsCount:      r.StarsCount,
		OpenIssuesCount: r.OpenIssuesCount,
		WatchersCount:   r.WatchersCount,
		StartDate:       r.StartDate.String(),
		EndDate:         r.EndDate.String(),
	}
}
