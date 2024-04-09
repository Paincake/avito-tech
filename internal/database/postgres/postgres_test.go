package postgres

import (
	"github.com/Paincake/avito-tech/internal/dto"
	"testing"
)

func TestKal(t *testing.T) {
	banner := dto.Banner{
		Tags:      []int64{1, 2},
		FeatureId: 1,
		Content: dto.Content{
			Title: "",
			Text:  "",
			Url:   "",
		},
		IsActive: true,
	}
	db, _ := New("avito", "avito", "avito", "localhost", "5432")
	_, err := db.InsertBanner(banner)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
