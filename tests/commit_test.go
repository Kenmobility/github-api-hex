package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SaveCommit(t *testing.T) {
	cr := repositories.NewGormCommitRepository(db)
	rr := repositories.NewGormRepositoryRepository(db)

	repo := createRandomRepository(t, rr)
	_ = createRandomCommit(t, repo, cr)
}

func Test_AllCommitsByRepository(t *testing.T) {
	query := dtos.APIPagingDto{}

	cr := repositories.NewGormCommitRepository(db)
	rr := repositories.NewGormRepositoryRepository(db)

	repo := createRandomRepository(t, rr)

	_ = createRandomCommit(t, repo, cr)
	_ = createRandomCommit(t, repo, cr)
	_ = createRandomCommit(t, repo, cr)

	commitRes, err := cr.AllCommitsByRepository(context.Background(), repo, query)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(commitRes.Commits), 3)

	require.Equal(t, repo.Name, commitRes.Commits[0].Repository)
	require.Equal(t, repo.Name, commitRes.Commits[1].Repository)
	require.Equal(t, repo.Name, commitRes.Commits[2].Repository)
}

func createRandomCommit(t *testing.T, repo domain.Repository, ci domain.CommitRepository) domain.Commit {
	require.NotEmpty(t, repo)
	require.NotZero(t, repo.ID)

	commit := domain.Commit{
		CommitID:     helpers.RandomString(12),
		RepositoryID: repo.ID,
		Message:      helpers.RandomWords(6),
		Author:       helpers.RandomString(8),
		Date:         time.Now(),
		URL:          fmt.Sprintf("%s/commit/%s", repo.URL, helpers.RandomString(20)),
	}

	nCommit, err := ci.SaveCommit(context.Background(), commit)
	require.NoError(t, err)
	require.NotEmpty(t, nCommit)

	require.Equal(t, repo.ID, nCommit.RepositoryID)
	require.Equal(t, nCommit.CommitID, commit.CommitID)
	require.Equal(t, nCommit.Message, commit.Message)
	require.Equal(t, nCommit.Author, commit.Author)
	require.Equal(t, nCommit.URL, commit.URL)

	return *nCommit
}
