package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateRepository(t *testing.T) {
	r := repositories.NewGormRepositoryRepository(db)
	createRandomRepository(t, r)
}

func Test_GetRepository(t *testing.T) {
	r := repositories.NewGormRepositoryRepository(db)

	repo := createRandomRepository(t, r)

	gotRepo, err := r.RepositoryByPublicId(context.Background(), repo.PublicID)
	require.NoError(t, err)
	require.NotEmpty(t, gotRepo)

	require.Equal(t, repo.PublicID, gotRepo.PublicID)
	require.Equal(t, repo.Name, gotRepo.Name)
	require.Equal(t, repo.URL, gotRepo.URL)

	require.WithinDuration(t, repo.CreatedAt, gotRepo.CreatedAt, time.Second)
}

func Test_GetAllRepositories(t *testing.T) {
	r := repositories.NewGormRepositoryRepository(db)

	_ = createRandomRepository(t, r)
	_ = createRandomRepository(t, r)
	_ = createRandomRepository(t, r)

	repos, err := r.AllRepositories(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(repos), 3)
}

func Test_TrackRepository(t *testing.T) {
	r := repositories.NewGormRepositoryRepository(db)

	newRepo := createRandomRepository(t, r)

	_, err := r.TrackRepository(context.Background(), newRepo)
	require.NoError(t, err)

	// fetch the created repo
	fRepo, err := r.RepositoryByPublicId(context.Background(), newRepo.PublicID)
	require.NoError(t, err)
	require.Equal(t, fRepo.IsTracking, true)

	tRepo, err := r.TrackedRepository(context.Background())
	require.NoError(t, err)
	require.Equal(t, fRepo, tRepo)
}

func Test_AtMostOneRepoTracking(t *testing.T) {
	r := repositories.NewGormRepositoryRepository(db)

	_ = createRandomRepository(t, r)
	repo2 := createRandomRepository(t, r)
	_ = createRandomRepository(t, r)

	//ensure repo2 is not empty
	require.NotEmpty(t, repo2)

	// ensure we have atleast 3 created repositories
	repos, err := r.AllRepositories(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(repos), 3)

	// set repo2 to be tracked
	tRepo, err := r.TrackRepository(context.Background(), repo2)
	require.NoError(t, err)
	require.NotEmpty(t, tRepo)

	trackedRepo, err := r.TrackedRepository(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, trackedRepo)

	// fetch the repo
	fRepo, err := r.RepositoryByPublicId(context.Background(), repo2.PublicID)
	require.NoError(t, err)
	require.Equal(t, fRepo.IsTracking, true)

	// ensure that the tracked repository is same as repo set for tracking
	require.Equal(t, trackedRepo, fRepo)
}

func createRandomRepository(t *testing.T, ri domain.RepositoryRepository) domain.Repository {
	arg := domain.Repository{
		PublicID:    uuid.New().String(),
		Name:        helpers.RandomRepositoryName(),
		Description: "",
		URL:         helpers.RandomRepositoryUrl(),
		StartDate:   helpers.RandomFetchStartDate(),
		EndDate:     helpers.RandomFetchEndDate(),
	}

	repo, err := ri.SaveRepository(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, repo)
	require.NotEmpty(t, repo.CreatedAt)
	require.NotEmpty(t, repo.StartDate)
	require.NotEmpty(t, repo.EndDate)

	require.Empty(t, repo.Description)

	require.Equal(t, arg.Name, repo.Name)
	require.Equal(t, arg.PublicID, repo.PublicID)
	require.Equal(t, repo.IsTracking, false)

	require.NotZero(t, repo.ID)

	return *repo
}
