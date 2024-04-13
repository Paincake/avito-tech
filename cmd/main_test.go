package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Paincake/avito-tech/internal/config"
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/database/postgres"
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/Paincake/avito-tech/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	TableCreationDDL = `
	CREATE TABLE IF NOT EXISTS api_users (
		username varchar PRIMARY KEY,
		password varchar,
		role varchar
	);
CREATE TABLE IF NOT EXISTS features (
    feature_id serial PRIMARY KEY ,
    description text
);

CREATE TABLE IF NOT EXISTS tags (
    tag_id serial PRIMARY KEY ,
    description text
);

CREATE TABLE IF NOT EXISTS banners(
    banner_id serial PRIMARY KEY ,
    feature_id int REFERENCES features(feature_id) ON DELETE CASCADE,
    content_title text,
    content_text text,
    content_url text,
    is_active bool,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS banner_tags (
    banner_id int REFERENCES banners(banner_id),
    tag_id int REFERENCES tags(tag_id),
    PRIMARY KEY (banner_id, tag_id)
)
`
	TableDeletionDDL = `
	TRUNCATE TABLE banner_tags CASCADE;
	TRUNCATE TABLE banners CASCADE;
	TRUNCATE TABLE features CASCADE;
	TRUNCATE TABLE tags CASCADE;
	ALTER SEQUENCE banners_banner_id_seq RESTART;
	ALTER SEQUENCE features_feature_id_seq RESTART;
	ALTER SEQUENCE tags_tag_id_seq RESTART;
`
	TableFillDDL = `
INSERT INTO features (description) VALUES ('f1'), ('f2'), ('f3');
INSERT INTO tags (description) VALUES ('t1'), ('t2'), ('t4');

INSERT INTO banners (feature_id, content_title, content_text, content_url, is_active, created_at, updated_at)
VALUES
(1, 'a', 'b', 'c', true, '2024-04-12 09:23:51.447097 +00:00'::timestamptz,  '2024-04-12 09:23:51.447097 +00:00'::timestamptz),
(2, 'a', 'b', 'c', true, '2024-04-12 09:23:51.447097 +00:00'::timestamptz,  '2024-04-12 09:23:51.447097 +00:00'::timestamptz),
(3, 'a', 'b', 'c', false,'2024-04-12 09:23:51.447097 +00:00'::timestamptz,  '2024-04-12 09:23:51.447097 +00:00'::timestamptz);

INSERT INTO banner_tags
VALUES 
(1, 1),
(1, 2),
(2, 2),
(2, 3),
(3, 3);
`
)

var router *echo.Echo
var db database.BannerRepository
var done chan bool

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func TestServerLoad(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/user_banner?tag_id=%d&feature_id=%d", rand.IntN(3), rand.IntN(3)), nil)
			token, _ := server.CreateJWT("admin")
			req.Header.Set("Token", token)
			router.ServeHTTP(recorder, req)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestPostBanner_ShouldGet201(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Banner{
		Tags:      []int64{1},
		FeatureId: 1,
		Content:   dto.Content{"a", "v", "c"},
		IsActive:  false,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(body))
	token, _ := server.CreateJWT("admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Result().StatusCode)
}

func TestPostBanner_ShouldThrow403(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Banner{
		Tags:      []int64{1},
		FeatureId: 1,
		Content:   dto.Content{"a", "v", "c"},
		IsActive:  false,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(body))
	token, _ := server.CreateJWT("user")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Result().StatusCode)
}

func TestGetUserBannerWithoutParams_ShouldThrow400(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user_banner", nil)
	token, _ := server.CreateJWT("admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
}

func TestGetBannersWithoutParams_ShouldReturnGivenValues(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/banner", nil)
	token, _ := server.CreateJWT("admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	ti, _ := time.Parse(time.RFC3339, "2024-04-12T09:23:51.447097Z")
	examples := []dto.Banner{
		{
			[]int64{1, 2},
			1,
			dto.Content{
				Title: "a", Text: "b", Url: "c",
			},
			true,
			ti.Format(time.RFC3339),
			ti.Format(time.RFC3339),
		},
		{
			[]int64{2, 3},
			2,
			dto.Content{
				Title: "a", Text: "b", Url: "c",
			},
			true,
			ti.Format(time.RFC3339),
			ti.Format(time.RFC3339),
		},
		{
			[]int64{3},
			3,
			dto.Content{
				Title: "a", Text: "b", Url: "c",
			},
			false,
			ti.Format(time.RFC3339),
			ti.Format(time.RFC3339),
		},
	}

	decoder := json.NewDecoder(recorder.Result().Body)
	var banners []dto.Banner
	decoder.Decode(&banners)
	for i := range banners {
		examples[i].CreatedAt = banners[i].CreatedAt
		examples[i].UpdatedAt = banners[i].UpdatedAt
	}
	if !reflect.DeepEqual(banners, examples) {
		t.Fail()
	}
}

func TestGetBannersWithParams_ShouldReturnGivenValues(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/banner?tag_id=2&feature_id=2&limit=1&offset=0", nil)
	token, _ := server.CreateJWT("admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	ti, _ := time.Parse(time.RFC3339, "2024-04-12T09:23:51.447097Z")
	examples := []dto.Banner{
		{
			[]int64{2, 3},
			2,
			dto.Content{
				Title: "a", Text: "b", Url: "c",
			},
			true,
			ti.Format(time.RFC3339),
			ti.Format(time.RFC3339),
		},
	}

	decoder := json.NewDecoder(recorder.Result().Body)
	var banners []dto.Banner
	decoder.Decode(&banners)
	for i := range banners {
		examples[i].CreatedAt = banners[i].CreatedAt
		examples[i].UpdatedAt = banners[i].UpdatedAt
	}
	if !reflect.DeepEqual(banners, examples) {
		t.Fail()
	}
}

func setup() {
	var err error
	configPath := os.Getenv("TEST_CONFIG_PATH")
	if configPath == "" {
		log.Fatal(err)
	}
	cfg, err := config.MustLoad(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db, err = postgres.New("test_database", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	if err != nil {
		panic(err)
	}
	err = db.RunMigrations(TableCreationDDL, TableFillDDL)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	done = make(chan bool)
	cache := server.NewMemoryCache(db, 2.5, 1, done)
	ConfigureServer(db, cache, e, server.VerifyJWT)
	router = e

}

func teardown() {
	err := db.RunMigrations(TableDeletionDDL)
	if err != nil {
		panic(err)
	}
	close(done)
}
