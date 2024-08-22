package unit

import (
	"testing"
	"time"

	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormCommitRepository_SaveCommit(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&domain.Commit{}, &domain.Repository{})

	repo := repositories.NewGormCommitRepository(db)

	commit := &domain.Commit{
		ID:      "afa-afaf-afdaf",
		Message: "Initial commit",
		Author:  "Author1",
		Date:    time.Now(),
		URL:     "http://example.com",
	}

	err := repo.SaveCommit(commit)

	assert.Nil(t, err)

	var savedCommit domain.Commit
	db.First(&savedCommit, "sha = ?", "test-sha")

	assert.Equal(t, "test-sha", savedCommit.SHA)
	assert.Equal(t, "Initial commit", savedCommit.Message)
}
