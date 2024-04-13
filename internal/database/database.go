package database

import (
	"github.com/Paincake/avito-tech/internal/dto"
	"strconv"
	"strings"
	"time"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type BannerRepository interface {
	InsertBanner(banner dto.Banner) (int64, error)
	UpdateBannerById(id int64, banner dto.Banner) error
	DeleteBannerById(id int64) error
	SelectUserBanner(params dto.GetUserBannerParams) (UserBanner, error)
	SelectBanners(params dto.GetBannerParams) ([]Banner, error)
	Login(username string, password string) (string, error)
	Signup(username string, password string) error
	RunMigrations(query ...string) error
}

type Banner struct {
	BannerID     int64     `db:"banner_id"`
	TagIDs       string    `db:"tag_ids"`
	FeatureID    int64     `db:"feature_id"`
	ContentTitle string    `db:"content_title"`
	ContentText  string    `db:"content_text"`
	ContentURL   string    `db:"content_url"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func ConvertBannerToDto(banner Banner) (dto.Banner, error) {
	var ids []int64
	for _, id := range strings.Split(banner.TagIDs[1:len(banner.TagIDs)-1], ",") {
		validID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return dto.Banner{}, err
		}
		ids = append(ids, validID)
	}
	return dto.Banner{
		Tags:      ids,
		FeatureId: banner.FeatureID,
		Content: dto.Content{
			Title: banner.ContentTitle,
			Text:  banner.ContentText,
			Url:   banner.ContentURL,
		},
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt.Format(time.RFC3339),
		UpdatedAt: banner.UpdatedAt.Format(time.RFC3339),
	}, nil
}

type UserBanner struct {
	Title     string    `db:"content_title"`
	Text      string    `db:"content_text"`
	URL       string    `db:"content_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ConvertUserBannerToDto(banner UserBanner) dto.Content {
	return dto.Content{
		Title: banner.Title,
		Text:  banner.Text,
		Url:   banner.URL,
	}
}

type User struct {
	Username string `db:"username" required:"true"`
	Password string `db:"password" required:"true"`
	Role     string `db:"role" required:"true"`
}
