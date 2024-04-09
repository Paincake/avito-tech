package server

import (
	"encoding/json"
	"fmt"
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
	"net/http"
)

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetBanner converts echo context to params.
func (w *ServerInterfaceWrapper) GetBanner(ctx echo.Context) error {
	params, err := NewGetBannerParams(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter feature_id: %s", err))
	}
	err = w.Handler.GetBanner(ctx, *params)
	return err
}

// PostBanner converts echo context to params.
func (w *ServerInterfaceWrapper) PostBanner(ctx echo.Context) error {
	var banner dto.Banner
	body := ctx.Request().Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&banner)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body"))
	}
	err = validator.Validate(banner)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body"))
	}
	err = w.Handler.PostBanner(ctx, banner)
	return err
}

// DeleteBannerId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteBannerId(ctx echo.Context) error {
	params, err := NewDeleteBannerIdParams(ctx)
	err = w.Handler.DeleteBannerId(ctx, *params)
	return err
}

// PatchBannerId converts echo context to params.
func (w *ServerInterfaceWrapper) PatchBannerId(ctx echo.Context) error {
	var banner dto.Banner
	body := ctx.Request().Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&banner)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body"))
	}
	err = validator.Validate(banner)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request body"))
	}

	params, err := NewPatchBannerIdParams(ctx)
	err = w.Handler.PatchBannerId(ctx, *params)
	return err
}

// GetUserBanner converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserBanner(ctx echo.Context) error {
	params, err := NewGetUserBannerParams(ctx)
	err = w.Handler.GetUserBanner(ctx, *params)
	return err
}
