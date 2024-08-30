package controllers

import (
	"context"

	"github.com/kenmobility/github-api-hex/common/message"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
)

type CommitController interface {
	GetAllCommitsByRepository(ctx context.Context, repoId string, query dtos.APIPagingDto) (string, *dtos.AllCommitsResponse, error)
	GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (string, []string, error)
}

type commitController struct {
	commitRepo domain.CommitRepository
	repoRepo   domain.RepositoryRepository
}

func NewCommitController(commitRepo domain.CommitRepository, repoRepo domain.RepositoryRepository) CommitController {
	return &commitController{
		commitRepo: commitRepo,
		repoRepo:   repoRepo,
	}
}

func (c *commitController) GetAllCommitsByRepository(ctx context.Context, repoId string, query dtos.APIPagingDto) (string, *dtos.AllCommitsResponse, error) {
	repo, err := c.repoRepo.RepositoryByPublicId(ctx, repoId)
	if err != nil {
		return "", nil, message.ErrInvalidRepositoryId
	}

	commits, err := c.commitRepo.AllCommitsByRepository(ctx, *repo, query)
	if err != nil {
		return "", nil, err
	}
	return repo.Name, commits, nil
}

func (c *commitController) GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (string, []string, error) {
	repo, err := c.repoRepo.RepositoryByPublicId(ctx, repoId)
	if err != nil {
		return "", nil, message.ErrInvalidRepositoryId
	}
	authors, err := c.commitRepo.TopCommitAuthorsByRepository(ctx, *repo, limit)
	if err != nil {
		return "", nil, err
	}

	return repo.Name, authors, nil
}
