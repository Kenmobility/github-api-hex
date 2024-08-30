package dtos

import (
	"time"
)

type (
	GithubCommitResponse struct {
		SHA     string `json:"sha"`
		NodeId  string `json:"node_id"`
		Commit  Commit `json:"commit"`
		URL     string `json:"url"`
		HtmlURL string `json:"html_url"`
	}

	Commit struct {
		Author  Author `json:"author"`
		Message string `json:"message"`
		URL     string `json:"url"`
	}

	Author struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	}
)

type CommitResponse struct {
	CommitID   string    `json:"commit_id"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	Date       time.Time `json:"date"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AllCommitsResponse struct {
	Commits  []CommitResponse `json:"commits"`
	PageInfo PagingInfo       `json:"page_info"`
}
