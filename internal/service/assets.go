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
	var assets = []string{}
	var urls = []string{}
	timeKey := time.Now().Format("20060102150405")
	for _, file := range input.Files {
		key := "Warehouse@cloud/Shotwot Originals/" + input.BriefId.String() + input.UserId + timeKey + "/" + file.Name
		url, err := b.s3.PresignedUrl(key, file.Filetype)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url.URL)
		assets = append(assets, key)
	}

	input.AssetFile = assets
	asset, err := b.repo.Create(ctx, input.Asset)
	if err != nil {
		return nil, err
	}
	return &domain.AssetRes{
		Urls:  urls,
		Asset: asset,
	}, nil

}
