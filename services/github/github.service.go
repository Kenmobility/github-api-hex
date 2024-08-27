package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kenmobility/github-api-hex/common/client"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"github.com/kenmobility/github-api-hex/services"
)

type GitHubAPIClient struct {
	baseURL              string
	token                string
	fetchInterval        time.Duration
	commitRepository     domain.CommitRepository
	repositoryRepository domain.RepositoryRepository
	client               *client.RestClient
}

func (g *GitHubAPIClient) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", g.token),
	}
}

func NewGitHubAPIClient(baseUrl string, token string, fetchInterval time.Duration, commitRepository domain.CommitRepository, repositoryRepository domain.RepositoryRepository) services.RepositoryTracker {
	client := client.NewRestClient()

	gc := GitHubAPIClient{
		baseURL:              baseUrl,
		token:                token,
		fetchInterval:        fetchInterval,
		commitRepository:     commitRepository,
		repositoryRepository: repositoryRepository,
		client:               client,
	}
	ts := services.RepositoryTracker(&gc)
	return ts
}

func (g *GitHubAPIClient) FetchAndSaveCommits(ctx context.Context, repo domain.Repository, since time.Time, until time.Time) ([]domain.Commit, error) {
	var result []domain.Commit

	endpoint := fmt.Sprintf("%s/repos/%s/commits?since=%s&until=%s", g.baseURL, repo.Name, since.Format(time.RFC3339), until.Format(time.RFC3339))
	for endpoint != "" {
		//check if repo is still the tracked repo before calling for next page data
		trackedRepo, err := g.repositoryRepository.TrackedRepository(ctx)
		if err != nil {
			return nil, err
		}

		if trackedRepo.ID != repo.ID { //tracked repo has changed
			return nil, errors.New("tracked repo changed")
		}

		commitRes, nextURL, err := g.fetchCommitsPage(endpoint)
		if err != nil {
			return nil, err
		}

		for _, c := range commitRes {
			result = append(result, domain.Commit{
				CommitID: c.SHA,
				Message:  c.Commit.Message,
				Author:   c.Commit.Author.Name,
				Date:     c.Commit.Author.Date,
				URL:      c.HtmlURL,
			})
		}

		for _, commit := range result {
			commit.RepositoryID = repo.ID

			if err := g.commitRepository.SaveCommit(ctx, commit); err != nil {
				log.Printf("Error saving commit id-%s: %v\n", commit.CommitID, err)
			}
		}

		endpoint = nextURL
	}

	return result, nil
}

func (g GitHubAPIClient) StartTracking(fetchInterval time.Duration) {
	for {
		g.runRepositoryTracker()
		time.Sleep(fetchInterval)
	}
}

func (g GitHubAPIClient) runRepositoryTracker() {
	ctx := context.Background()

	trackedRepo, err := g.repositoryRepository.TrackedRepository(ctx)
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
	g.FetchAndSaveCommits(ctx, *trackedRepo, trackedRepo.StartDate, trackedRepo.EndDate)
}

func (g *GitHubAPIClient) fetchCommitsPage(url string) ([]dtos.GithubCommitResponse, string, error) {

	response, err := g.client.Get(url, map[string]string{}, g.getHeaders())
	if err != nil {
		log.Println("error fetching commits: ", err)
		return nil, "", nil
	}

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch commits; status code: %v", response.StatusCode)
	}

	var commitRes []dtos.GithubCommitResponse

	if err := json.Unmarshal([]byte(response.Body), &commitRes); err != nil {
		fmt.Printf("marshal error, [%v]", err)
		return nil, "", errors.New("could not unmarshal commits response")
	}

	nextURL := g.parseNextURL(response.Headers["Link"])

	return commitRes, nextURL, nil
}

func (api *GitHubAPIClient) parseNextURL(linkHeader []string) string {
	if len(linkHeader) == 0 {
		return ""
	}

	links := strings.Split(linkHeader[0], ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}

		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])

		if relPart == `rel="next"` {
			return urlPart
		}
	}

	return ""
}
