package postgres

import (
	"fmt"
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/dto"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type EntityNotFound struct {
	Err error
}

func (b EntityNotFound) Error() string {
	return b.Err.Error()
}

type OptionalTagIdParam struct {
	Value int64
}

type OptionalFeatureIdParam struct {
	Value int64
}

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

func (d *Database) RunMigrations(query ...string) error {
	for _, q := range query {
		_, err := d.db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func (d *Database) UpdateBannerById(id int64, banner dto.Banner) error {
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
		return EntityNotFound{Err: err}
	}
	if len(banner.Tags) > 0 {
		_, err := d.db.Query(`INSERT INTO banner_tags VALUES ($1, unnest($2::INTEGER[]))`, id, banner.Tags)
		if err != nil {
			return fmt.Errorf("error inserting banner tags: %s", err)
		}
	}
	return nil
}

func (d *Database) DeleteBannerById(id int64) error {
	result, err := d.db.Exec(`DELETE FROM banners WHERE banner_id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting banner: %s", err)
	}
	affected, err := result.RowsAffected()
	if affected == 0 {
		return EntityNotFound{err}
	}
	return nil
}

func (d *Database) SelectUserBanner(params dto.GetUserBannerParams) (database.UserBanner, error) {
	var banner database.UserBanner
	err := d.db.Get(&banner,
		`SELECT b.content_title, b.content_text, b.content_url, b.created_at, b.updated_at FROM banners b 
                    JOIN banner_tags bt ON bt.banner_id = b.banner_id
					WHERE b.feature_id= $1 AND bt.tag_id = $2
					AND b.is_active = (CASE WHEN $3 = true THEN true ELSE b.is_active END)
					`, params.FeatureId, params.TagId, params.UseActive)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return banner, EntityNotFound{Err: err}
		}
		return database.UserBanner{}, fmt.Errorf("error selecting user banner: %s", err)
	}
	empty := banner == database.UserBanner{}
	if empty {
		return banner, EntityNotFound{Err: err}
	}
	return banner, nil
}

func (d *Database) SelectBanners(params dto.GetBannerParams) ([]database.Banner, error) {
	var banners []database.Banner
	err := d.db.Select(&banners,
		`SELECT b.banner_id, tags.tag_ids, b.feature_id, b.content_title, b.content_text, b.content_url, b.is_active, b.created_at, b.updated_at FROM banners b
			   JOIN banner_tags bt ON bt.banner_id = b.banner_id

			   JOIN 
					(
						SELECT b.banner_id, array_agg(bt.tag_id) as tag_ids FROM banners b
					 	JOIN banner_tags bt on b.banner_id = bt.banner_id
              		 	GROUP BY b.banner_id, b.feature_id
						HAVING b.feature_id = (CASE WHEN $1 = $6::int THEN b.feature_id ELSE $1 END) 
					) tags ON tags.banner_id = b.banner_id

 			   WHERE
					bt.tag_id = (CASE WHEN $2 = $6::int THEN bt.tag_id ELSE $2 END) AND 
					b.is_active = (CASE WHEN $3 = true THEN true ELSE b.is_active END)
 			   
			   GROUP BY b.banner_id, tags.tag_ids
			   LIMIT $4 OFFSET $5`,
		params.FeatureId,
		params.TagId,
		params.UseActive,
		params.Limit,
		params.Offset,
		dto.DefaultIdValue,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting banners: %s", err)
	}
	return banners, nil
}
func (d *Database) Login(username string, password string) (string, error) {
	var user database.User
	err := d.db.Get(&user, "SELECT username, password, role FROM api_users WHERE username = $1", username)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return user.Role, nil
}

func (d *Database) Signup(username string, password string) error {
	_, err := d.db.Query("INSERT INTO api_users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}
