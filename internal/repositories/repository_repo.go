package repositories

import (
	"context"

	"github.com/kenmobility/github-api-hex/common/message"
	"github.com/kenmobility/github-api-hex/internal/domain"
	"gorm.io/gorm"
)

type GormRepositoryRepository struct {
	DB *gorm.DB
}

func NewGormRepositoryRepository(db *gorm.DB) domain.RepositoryRepository {
	return &GormRepositoryRepository{DB: db}
}

func (r *GormRepositoryRepository) SaveRepository(ctx context.Context, repository domain.Repository) (*domain.Repository, error) {
	err := r.DB.WithContext(ctx).Create(&repository).Error
	return &repository, err
}

func (r *GormRepositoryRepository) RepositoryByPublicId(ctx context.Context, publicId string) (*domain.Repository, error) {
	var repo domain.Repository
	err := r.DB.WithContext(ctx).Where("public_id = ?", publicId).Find(&repo).Error

	if repo.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return &repo, err
}

func (r *GormRepositoryRepository) AllRepositories(ctx context.Context) ([]domain.Repository, error) {
	var repositories []domain.Repository

	err := r.DB.WithContext(ctx).Find(&repositories).Error
	return repositories, err
}

func (r *GormRepositoryRepository) TrackedRepository(ctx context.Context) (*domain.Repository, error) {
	var repo domain.Repository
	err := r.DB.WithContext(ctx).Where("is_tracking = ?", true).First(&repo).Error

	if repo.ID == 0 {
		return nil, message.ErrNoTrackingRepositorySet
	}

	return &repo, err
}

func (r *GormRepositoryRepository) TrackRepository(ctx context.Context, repository domain.Repository) (*domain.Repository, error) {
	// reset all repositories to not tracking
	err := r.DB.WithContext(ctx).Model(&domain.Repository{}).Where("is_tracking = ?", true).Update("is_tracking", false).Error
	if err != nil {
		return nil, err
	}
	// Set the specified repository to tracking
	err = r.DB.WithContext(ctx).Model(&domain.Repository{}).Where("public_id = ?", repository.PublicID).
		Updates(&repository).Error

	return &repository, err
}
