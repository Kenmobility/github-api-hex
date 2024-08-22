package unit

import (
	"context"
	"testing"
	"time"

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
