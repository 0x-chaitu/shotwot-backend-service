package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/s3"
)

type BriefApplicationsService struct {
	repo repository.BriefApplications
	s3   *s3.S3Client
}

func NewBriefApplicationsService(repo repository.BriefApplications, s3 *s3.S3Client) *BriefApplicationsService {
	return &BriefApplicationsService{
		repo: repo,
		s3:   s3,
	}
}

func (b *BriefApplicationsService) Create(ctx context.Context, input domain.BriefApplication) (*domain.BriefApplication, error) {

	input.Opened = false
	input.Status = domain.Applied
	briefapplication, err := b.repo.Create(ctx, &input)
	if err != nil {
		return nil, err
	}
	return briefapplication, nil

}

func (b *BriefApplicationsService) GetBriefApplications(ctx context.Context, id string) ([]*domain.BriefApplication, error) {
	return b.repo.GetBriefApplications(ctx, id)
}

func (b *BriefApplicationsService) GetBriefApplication(ctx context.Context, id string) (*domain.UserBriefAppliedDetails, error) {
	return b.repo.GetBriefApplication(ctx, id)
}

func (b *BriefApplicationsService) UpdateBriefApplication(ctx context.Context, input *domain.BriefApplication) (*domain.BriefApplication, error) {
	return b.repo.Update(ctx, input)
}
