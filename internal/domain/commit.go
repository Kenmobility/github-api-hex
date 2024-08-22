package domain

import (
	"context"
	"time"

	"github.com/kenmobility/github-api-hex/dtos"
)

type Commit struct {
	ID           uint       `gorm:"primaryKey"`
	CommitID     string     `gorm:"type:varchar(100);uniqueIndex"`
	Message      string     `gorm:"type:varchar"`
	Author       string     `gorm:"type:varchar"`
	Date         time.Time  `gorm:"index"`
	URL          string     `gorm:"type:varchar"`
	RepositoryID uint       `gorm:"index"` //Foreign key to Repository
	Repository   Repository `gorm:"foreignKey:RepositoryID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CommitRepository interface {
	SaveCommit(ctx context.Context, commit Commit) error
	AllCommitsByRepository(ctx context.Context, repo Repository, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error)
	TopCommitAuthorsByRepository(ctx context.Context, repository Repository, limit int) ([]string, error)
}

func (c Commit) ToDto() dtos.CommitResponse {
	return dtos.CommitResponse{
		ID:         c.ID,
		CommitID:   c.CommitID,
		Message:    c.Message,
		Author:     c.Author,
		Date:       c.Date,
		URL:        c.URL,
		Repository: c.RepositoryName(),
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func (c *Commit) RepositoryName() string {
	if c == nil {
		return ""
	}
	return c.Repository.Name
}
