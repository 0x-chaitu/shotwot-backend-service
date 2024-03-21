package service

import (
	"context"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/repository"
	"shotwot_backend/pkg/s3"
	"time"
)

type AssetsService struct {
	repo repository.Assets
	s3   *s3.S3Client
}

func NewAssetsService(repo repository.Assets, s3 *s3.S3Client) *AssetsService {
	return &AssetsService{
		repo: repo,
		s3:   s3,
	}
}

func (b *AssetsService) Create(ctx context.Context, input *domain.AssetInput) (*domain.AssetRes, error) {
	timeKey := time.Now().Format("20060102150405")
	key := "Warehouse@cloud/Shotwot Originals/" + input.BriefId.String() + "/" + input.UserId + "/" + timeKey + "/" + input.File.Name
	url, err := b.s3.PresignedUrl(key, input.File.Filetype)
	if err != nil {
		return nil, err
	}
	input.Rating = &domain.Rating{
		Current: 0,
		Total:   0,
	}

	input.AssetFile = key
	asset, err := b.repo.Create(ctx, input.Asset)
	if err != nil {
		return nil, err
	}
	return &domain.AssetRes{
		Url:   url.URL,
		Asset: asset,
	}, nil
}

func (b *AssetsService) Update(ctx context.Context, input *domain.Asset) (*domain.Asset, error) {
	input.Updated = time.Now()
	asset, err := b.repo.Update(ctx, input)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (b *AssetsService) GetAllAssets(ctx context.Context) ([]*domain.Asset, error) {
	return b.repo.GetAssets(ctx)
}
