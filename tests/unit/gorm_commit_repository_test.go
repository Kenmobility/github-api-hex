package unit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestGormCommitRepository_SaveCommit(t *testing.T) {

	repo := repositories.NewGormCommitRepository(db)

	commit := &domain.Commit{
		CommitID: "afa-afaf-afdaf",
		Message:  "Initial commit",
		Author:   "Author1",
		Date:     time.Now(),
		URL:      "http://example.com",
	}

	err := repo.SaveCommit(context.Background(), *commit)

	assert.Nil(t, err)

	var savedCommit domain.Commit
	db.First(&savedCommit, "commit_id = ?", "afa-afaf-afdaf")

	assert.Equal(t, "Initial commit", savedCommit.Message)
	assert.Equal(t, "Author1", savedCommit.Author)
}

func TestGormCommitRepository_AllCommitsByRepository(t *testing.T) {
	cRepo := repositories.NewGormCommitRepository(db)

	var query = dtos.APIPagingDto{}

	newRepo := &domain.Repository{
		PublicID: uuid.New().String(),
		Name:     "owner/myrepo",
		URL:      "http://github.com/owner/myrepo",
	}

	db.Create(newRepo)

	commit := &domain.Commit{
		CommitID:     "test-id-test",
		RepositoryID: newRepo.ID,
		Message:      "Another awesome commit",
		Author:       "Author2",
		Date:         time.Now(),
		URL:          "http://githubexample.com",
	}

	db.Create(commit)

	commits, err := cRepo.AllCommitsByRepository(context.Background(), *newRepo, query)
	assert.Nil(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "test-id-test", commits.Commits[0].CommitID)
	assert.Equal(t, "Another awesome commit", commits.Commits[0].Message)
}
