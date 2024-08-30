package dtos

type AddRepositoryRequestDto struct {
	Name string `json:"name" validate:"required"`
}

type GitHubRepoMetadataResponse struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count"`
	StarsCount      int    `json:"stars_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	WatchersCount   int    `json:"watchers_count"`
}
