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

func (b *BriefApplicationsService) Create(ctx context.Context, input *domain.BriefApplicationInput) (*domain.BriefApplicationRes, error) {

	briefapplication, err := b.repo.Create(ctx, input.BriefApplication)
	if err != nil {
		return nil, err
	}
	return &domain.BriefApplicationRes{
		BriefApplication: briefapplication,
	}, nil

}
