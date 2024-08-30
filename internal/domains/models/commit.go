package models

import (
	"time"
)

type Commit struct {
	CommitID       string    `json:"commit_id"`
	Message        string    `json:"message"`
	Author         string    `json:"author"`
	Date           time.Time `json:"date"`
	URL            string    `json:"url"`
	RepositoryName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
