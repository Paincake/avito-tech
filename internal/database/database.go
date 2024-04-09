package database

import (
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/Paincake/avito-tech/internal/server"
)

type BannerRepository interface {
	InsertBanner(banner dto.Banner) (int64, error)
	UpdateBannerById(id int64, banner dto.Banner) error
	DeleteBannerById(id int64) error
	SelectUserBanner(params server.GetUserBannerParams) (UserBanner, error)
	SelectBanners(params server.GetBannerParams) ([]Banner, error)
}

type Banner struct {
	BannerId     int64   `db:"banner_id"`
	TagIds       []int64 `db:"tag_ids"`
	FeatureId    int64   `db:"feature_id"`
	ContentTitle string  `db:"title"`
	ContentText  string  `db:"text"`
	ContentUrl   string  `db:"url"`
	IsActive     bool    `db:"is_active"`
	CreatedAt    string  `db:"created_at"`
	UpdatedAt    string  `db:"updated_at"`
}

type UserBanner struct {
	Title string `db:"title"`
	Text  string `db:"text"`
	Url   string `db:"text"`
}
