package services

import (
	"context"
	"time"

	"github.com/kenmobility/github-api-hex/internal/domain"
)

type RepositoryTracker interface {
	FetchAndSaveCommits(ctx context.Context, repo domain.Repository, since time.Time, until time.Time) ([]domain.Commit, error)
	StartTracking(fetchInterval time.Duration)
}
