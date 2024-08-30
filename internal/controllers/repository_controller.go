package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/common/message"
	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
)

type RepositoryController interface {
	AddRepository(ctx context.Context, data dtos.AddRepositoryRequestDto) (*dtos.RepositoryResponse, error)
	TrackRepository(ctx context.Context, data dtos.TrackRepositoryRequestDto) (*dtos.RepositoryResponse, error)
	GetRepositoryById(ctx context.Context, id string) (*dtos.RepositoryResponse, error)
	GetAllRepositories(ctx context.Context) ([]dtos.RepositoryResponse, error)
}

type repositoryController struct {
	repositoryRepo domain.RepositoryRepository
	config         *config.Config
}

func NewRepositoryController(repositoryRepo domain.RepositoryRepository, config *config.Config) RepositoryController {
	return &repositoryController{
		repositoryRepo: repositoryRepo,
		config:         config,
	}
}

func (r *repositoryController) AddRepository(ctx context.Context, data dtos.AddRepositoryRequestDto) (*dtos.RepositoryResponse, error) {
	//validate repository name to ensure it has owner and repo name
	if !helpers.IsRepositoryNameValid(data.Name) {
		return nil, message.ErrInvalidRepositoryName
	}

	repository := domain.Repository{
		PublicID:    uuid.New().String(),
		Name:        data.Name,
		Description: data.Description,
		URL:         data.URL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo, err := r.repositoryRepo.SaveRepository(ctx, repository)
	if err != nil {
		return nil, err
	}

	repoResponse := repo.ToDto()

	return &repoResponse, nil
}

func (r *repositoryController) TrackRepository(ctx context.Context, data dtos.TrackRepositoryRequestDto) (*dtos.RepositoryResponse, error) {
	repo, err := r.repositoryRepo.RepositoryByPublicId(ctx, data.RepoId)

	if err != nil && err == message.ErrNoRecordFound {
		return nil, message.ErrRepositoryNotFound
	}

	if err != nil && err != message.ErrNoRecordFound {
		return nil, err
	}

	var startDate, endDate time.Time

	if data.StartDate != "" {
		startDate, err = time.Parse(time.RFC3339, data.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
	} else {
		startDate = r.config.DefaultStartDate
	}

	if data.EndDate != "" {
		endDate, err = time.Parse(time.RFC3339, data.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
	} else {
		endDate = r.config.DefaultEndDate
	}

	repo.StartDate = startDate
	repo.EndDate = endDate

	trackedRepo, err := r.repositoryRepo.TrackRepository(ctx, *repo)
	if err != nil {
		return nil, err
	}

	repoResponse := trackedRepo.ToDto()

	return &repoResponse, nil
}

func (r *repositoryController) GetRepositoryById(ctx context.Context, id string) (*dtos.RepositoryResponse, error) {
	repo, err := r.repositoryRepo.RepositoryByPublicId(ctx, id)

	if err != nil && err == message.ErrNoRecordFound {
		return nil, message.ErrRepositoryNotFound
	}

	if err != nil && err != message.ErrNoRecordFound {
		return nil, err
	}

	repoResponse := repo.ToDto()

	return &repoResponse, nil
}

func (r *repositoryController) GetAllRepositories(ctx context.Context) ([]dtos.RepositoryResponse, error) {
	repos, err := r.repositoryRepo.AllRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repositoryResponse := make([]dtos.RepositoryResponse, 0, len(repos))

	for _, repo := range repos {
		rr := dtos.RepositoryResponse{
			Id:              repo.PublicID,
			Name:            repo.Name,
			Description:     repo.Description,
			URL:             repo.URL,
			Language:        repo.Language,
			ForksCount:      repo.ForksCount,
			StarsCount:      repo.StarsCount,
			OpenIssuesCount: repo.OpenIssuesCount,
			WatchersCount:   repo.WatchersCount,
			StartDate:       repo.StartDate.String(),
			EndDate:         repo.EndDate.String(),
			IsTracking:      repo.IsTracking,
		}

		repositoryResponse = append(repositoryResponse, rr)
	}

	return repositoryResponse, nil
}
