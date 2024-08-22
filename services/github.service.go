package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/internal/integrations/api"
)

type GitHubService struct {
	api            *api.GitHubAPIClient
	commitRepo     domain.CommitRepository
	repositoryRepo domain.RepositoryRepository
	config         *config.Config
}

func NewGithubService(api *api.GitHubAPIClient, commitRepo domain.CommitRepository,
	repositoryRepo domain.RepositoryRepository, config *config.Config) *GitHubService {

	return &GitHubService{
		api:            api,
		commitRepo:     commitRepo,
		repositoryRepo: repositoryRepo,
		config:         config,
	}
}

func (s *GitHubService) run() {
	ctx := context.Background()

	trackedRepo, err := s.repositoryRepo.TrackedRepository(ctx)
	if err != nil {
		log.Printf("Error fetching tracked repository: %v", err)
		return
	}

	if trackedRepo == nil {
		log.Println("no repository set to track")
		return
	}
	fmt.Printf("********Github repository tracking started for repo %s************\n",
		trackedRepo.Name)
	s.fetchAndSaveCommits(ctx, *trackedRepo)
}

func (s *GitHubService) StartTracking() {
	go func() {
		for {
			s.run()
			time.Sleep(s.config.FetchInterval)
		}
	}()
}

func (s *GitHubService) fetchAndSaveCommits(ctx context.Context, repo domain.Repository) {
	_, err := s.api.FetchAndSaveCommits(ctx, repo, repo.StartDate, repo.EndDate)
	if err != nil {
		log.Printf("Error fetching commits for repository %s: %v", repo.Name, err)
		return
	}
}
