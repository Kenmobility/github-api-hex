package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"gorm.io/gorm"
)

type GormCommitRepository struct {
	DB *gorm.DB
}

func NewGormCommitRepository(db *gorm.DB) domain.CommitRepository {
	return &GormCommitRepository{DB: db}
}

// SaveCommit stores a repository commit into the database
func (gc *GormCommitRepository) SaveCommit(ctx context.Context, commit domain.Commit) error {
	return gc.DB.WithContext(ctx).Create(&commit).Error
}

// GetAllCommitsByRepositoryName fetches all stores commits by repository name
func (gc *GormCommitRepository) AllCommitsByRepository(ctx context.Context, repository domain.Repository, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error) {
	var commits []domain.Commit

	var count, queryCount int64

	queryInfo, offset := getPaginationInfo(query)

	//db := c.db.Db.WithContext(ctx).Joins("Repository").Where("repositories.id = ?", repo.ID) //.Find(&commits).Error
	db := gc.DB.WithContext(ctx).Model(&domain.Commit{}).Where(&domain.Commit{RepositoryID: repository.ID}).
		Preload("Repository")

	db.Count(&count)

	db = db.Offset(offset).Limit(queryInfo.Limit).
		Order(fmt.Sprintf("commits.%s %s", queryInfo.Sort, queryInfo.Direction)).
		Find(&commits)
	db.Count(&queryCount)

	if db.Error != nil {
		log.Println("fetch commits error", db.Error.Error())
		return nil, db.Error
	}

	pagingInfo := getPagingInfo(queryInfo, int(count))
	pagingInfo.Count = len(commits)

	return &dtos.AllCommitsResponse{
		Commits:  commitResponse(commits),
		PageInfo: pagingInfo,
	}, nil
}

func (gc *GormCommitRepository) TopCommitAuthorsByRepository(ctx context.Context, repository domain.Repository, limit int) ([]string, error) {
	var authors []string
	err := gc.DB.WithContext(ctx).Model(&domain.Commit{}).
		Select("author").
		Where("repository_id = ?", repository.ID).
		Group("author").
		Order("count(author) DESC").
		Limit(limit).
		Find(&authors).Error

	return authors, err
}

func commitResponse(commits []domain.Commit) []dtos.CommitResponse {
	if len(commits) == 0 {
		return nil
	}

	commitsResponse := make([]dtos.CommitResponse, 0, len(commits))

	for _, c := range commits {
		cr := dtos.CommitResponse{
			ID:         c.ID,
			CommitID:   c.CommitID,
			Message:    c.Message,
			Author:     c.Author,
			Date:       c.Date,
			URL:        c.URL,
			Repository: c.RepositoryName(),
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		}

		commitsResponse = append(commitsResponse, cr)
	}

	return commitsResponse
}
