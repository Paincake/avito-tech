package postgres

import (
	"fmt"
	"github.com/Paincake/avito-tech/internal/dto"
	"testing"
)

func TestNew_ShouldReturnValidConnection(t *testing.T) {
	_, err := New("avito", "avito", "avito", "localhost", "5432")
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestNew_ShouldReturnError(t *testing.T) {
	_, err := New("av", "it", "ato", "localost", "542")
	if err == nil {
		t.Fatalf("Invalid connection created")
	}
}

func TestDatabase_SelectBanners(t *testing.T) {
	params := dto.GetBannerParams{
		FeatureId: -1, TagId: -1, Limit: 1,
	}
	db, _ := New("avito", "avito", "avito", "localhost", "5432")
	banners, err := db.SelectBanners(params)
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Printf("%v\n", banners)
}

func TestDatabase_SelectUserBanner(t *testing.T) {
	params := dto.GetUserBannerParams{
		TagId:        1,
		FeatureId:    1,
		LastRevision: false,
		UseActive:    false,
	}

	db, _ := New("avito", "avito", "avito", "localhost", "5432")
	res, err := db.SelectUserBanner(params)
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Printf("%v\n", res)
}
