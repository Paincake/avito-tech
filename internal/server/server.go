package server

import (
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получение всех баннеров c фильтрацией по фиче и/или тегу
	// (GET /banner)
	GetBanner(ctx echo.Context, params GetBannerParams) error
	// Создание нового баннера
	// (POST /banner)
	PostBanner(ctx echo.Context, banner dto.Banner) error
	// Удаление баннера по идентификатору
	// (DELETE /banner/{id})
	DeleteBannerId(ctx echo.Context, params DeleteBannerIdParams) error
	// Обновление содержимого баннера
	// (PATCH /banner/{id})
	PatchBannerId(ctx echo.Context, params PatchBannerIdParams) error
	// Получение баннера для пользователя
	// (GET /user_banner)
	GetUserBanner(ctx echo.Context, params GetUserBannerParams) error
}
