package dtos

type AddRepositoryRequestDto struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"required"`
}

type TrackRepositoryRequestDto struct {
	RepoId    string `json:"repo_id" validate:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type RepositoryResponse struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count"`
	StarsCount      int    `json:"stars_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	WatchersCount   int    `json:"watchers_count"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	IsTracking      bool   `json:"is_tracking"`
}
