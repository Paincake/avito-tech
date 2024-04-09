package postgres

import (
	"fmt"
	"github.com/Paincake/avito-tech/internal/dto"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type Database struct {
	db *sqlx.DB
}

func New(dbname, username, password, host, port string) (*Database, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?&sslmode=disable",
		username,
		password,
		host,
		port,
		dbname)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Database{db: db}, nil
}

func (d *Database) RunMigrations(query ...string) {
	for _, q := range query {
		d.db.Query(q)
	}
}

// InsertBanner TODO transaction?
func (d *Database) InsertBanner(banner dto.Banner) (int64, error) {
	var lastInserted int64
	err := d.db.Get(&lastInserted,
		`INSERT INTO banners (feature_id, content_title, content_text, content_url, is_active, created_at, updated_at)
			   VALUES
			   ($1, $2, $3, $4, $5, $6, $7)
				RETURNING banner_id`,
		banner.FeatureId,
		banner.Content.Title,
		banner.Content.Text,
		banner.Content.Url,
		banner.IsActive,
		time.Now(),
		time.Now())
	if err != nil {
		return -1, fmt.Errorf("error inserting a banner: %s", err)
	}
	_, err = d.db.Query(`INSERT INTO banner_tags VALUES ($1, unnest($2::INTEGER[]))`, lastInserted, banner.Tags)
	if err != nil {
		return -1, fmt.Errorf("error inserting banner tags: %s", err)
	}
	return lastInserted, nil
}

// UpdateBannerById TODO how to throw 404 if banner is not found?
func (d *Database) UpdateBannerById(id int64, banner dto.Banner) error {
	if len(banner.Tags) > 0 {
		_, err := d.db.Query(`INSERT INTO banner_tags VALUES ($1, unnest($2::INTEGER[]))`, id, banner.Tags)
		if err != nil {
			return fmt.Errorf("error inserting banner tags: %s", err)
		}
	}
	result, err := d.db.Exec(
		`UPDATE banners SET feature_id = $1, content_title = $2, content_text=$3,content_url=$4,is_active=$5,updated_at=$6 WHERE banner_id = $7`,
		banner.FeatureId,
		banner.Content.Title,
		banner.Content.Text,
		banner.Content.Url,
		banner.IsActive,
		time.Now(),
		id)
	if err != nil {
		return fmt.Errorf("error updating banner: %s", err)
	}
	if affected, err := result.RowsAffected(); affected == 0 {
		return fmt.Errorf("banner not found: %s", err)
	}
	return nil
}

func (d *Database) DeleteBannerById(id int64) error {
	_, err := d.db.Query(`DELETE FROM banners WHERE banner_id = $1 CASCADE`, id)
	if err != nil {
		return fmt.Errorf("error deleting banner: %s", err)
	}
	return nil
}

//func (d *Database) SelectUserBanner(params server.GetUserBannerParams) (database.UserBanner, error) {
//
//}
//
//func (d *Database) SelectBanners(params server.GetBannerParams) ([]database.Banner, error) {
//
//}
