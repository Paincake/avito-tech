package server

import (
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/dto"
)

type ServerInterface interface {
	GetBanner(params dto.GetBannerParams) ([]dto.Banner, error)
	// PostBanner Создание нового баннера
	// (POST /banner)
	PostBanner(banner dto.Banner) (int64, error)
	// DeleteBannerID Удаление баннера по идентификатору
	// (DELETE /banner/{id})
	DeleteBannerID(params dto.DeleteBannerIdParams) error
	// PatchBannerID Обновление содержимого баннера
	// (PATCH /banner/{id})
	PatchBannerID(params dto.PatchBannerIdParams, banner dto.Banner) error
	// GetUserBanner Получение баннера для пользователя
	// (GET /user_banner)
	GetUserBanner(params dto.GetUserBannerParams) (dto.Content, error)
	Login(username, password string) (string, error)
	Signup(username, password string) error
}

type Server struct {
	Repository database.BannerRepository
	Cache      BannerCache
}

func (s *Server) GetBanner(params dto.GetBannerParams) ([]dto.Banner, error) {
	dbBanners, err := s.Repository.SelectBanners(params)
	if err != nil {
		return nil, err
	}
	var dtoBanners []dto.Banner
	for _, banner := range dbBanners {
		dtoBanner, err := database.ConvertBannerToDto(banner)
		if err != nil {
			return nil, err
		}
		dtoBanners = append(dtoBanners, dtoBanner)
	}
	return dtoBanners, nil

}

func (s *Server) PostBanner(banner dto.Banner) (int64, error) {
	id, err := s.Repository.InsertBanner(banner)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *Server) DeleteBannerID(params dto.DeleteBannerIdParams) error {
	err := s.Repository.DeleteBannerById(params.BannerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) PatchBannerID(params dto.PatchBannerIdParams, banner dto.Banner) error {
	err := s.Repository.UpdateBannerById(params.BannerId, banner)
	if err != nil {
		return err
	}
	return nil
}
func (s *Server) GetUserBanner(params dto.GetUserBannerParams) (dto.Content, error) {
	var banner database.UserBanner
	var err error
	if params.LastRevision {
		banner, err = s.Repository.SelectUserBanner(params)
	} else {
		banner, err = s.Cache.GetBanner(params.FeatureId, params)
	}
	if err != nil {
		return dto.Content{}, err
	}
	return database.ConvertUserBannerToDto(banner), nil
}
func (s *Server) Login(username, password string) (string, error) {
	role, err := s.Repository.Login(username, password)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (s *Server) Signup(username, password string) error {
	return s.Repository.Signup(username, password)
}
