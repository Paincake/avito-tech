package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Paincake/avito-tech/internal/config"
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/database/postgres"
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	"net/http"
	"strings"
)

type ServerInterfaceWrapper struct {
	Handler ServerInterface
	Options config.Config
}

func (w *ServerInterfaceWrapper) GetBanner(ctx echo.Context) error {
	role := ctx.Get(dto.TokenRoleContextKey)
	if role != database.AdminRole {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Forbidden"))
	}
	var err error
	params, err := dto.NewGetBannerParams(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter feature_id: %s", err))
	}
	banners := make([]dto.Banner, 0)
	banners, err = w.Handler.GetBanner(*params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("internal server error: %s", err))
	}

	return ctx.JSON(http.StatusOK, banners)
}

// PostBanner converts echo context to params.
func (w *ServerInterfaceWrapper) PostBanner(ctx echo.Context) error {
	role := ctx.Get(dto.TokenRoleContextKey)
	if role != database.AdminRole {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Forbidden"))
	}
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
	id, err := w.Handler.PostBanner(banner)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("internal server error: %s", err))
	}
	return ctx.JSON(http.StatusCreated, struct {
		BannerID int64 `json:"banner_id"`
	}{BannerID: id})
}

// DeleteBannerID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteBannerID(ctx echo.Context) error {
	role := ctx.Get(dto.TokenRoleContextKey)
	if role != database.AdminRole {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Forbidden"))
	}
	params, err := dto.NewDeleteBannerIdParams(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request parameter: %s", err))
	}
	err = w.Handler.DeleteBannerID(*params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("internal server error: %s", err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

// PatchBannerID converts echo context to params.
func (w *ServerInterfaceWrapper) PatchBannerID(ctx echo.Context) error {
	role := ctx.Get(dto.TokenRoleContextKey)
	if role != database.AdminRole {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Forbidden"))
	}
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

	params, err := dto.NewPatchBannerIdParams(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request parameter: %s", err))
	}
	err = w.Handler.PatchBannerID(*params, banner)
	if err != nil {
		var entityErr postgres.EntityNotFound
		ok := errors.As(err, &entityErr)
		if ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no banner for given tag"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("internal server error: %s", err))
	}
	return ctx.NoContent(http.StatusOK)
}

// GetUserBanner converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserBanner(ctx echo.Context) error {
	role := ctx.Get(dto.TokenRoleContextKey)
	if role != database.AdminRole && role != database.UserRole {
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Forbidden"))
	}
	params, err := dto.NewGetUserBannerParams(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid request parameter: %s", err))
	}
	banner, err := w.Handler.GetUserBanner(*params)
	if err != nil {
		var entityErr postgres.EntityNotFound
		ok := errors.As(err, &entityErr)
		if ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no banner found"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("internal server error: %s", err))
	}
	return ctx.JSON(http.StatusOK, banner)
}

func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	creds := ctx.Request().Header.Get("Authorization")
	if creds == "" || len(strings.Split(creds, " ")) < 2 {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Required Authorization header missing"))
	}
	creds = strings.Split(creds, " ")[1]
	raw, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Auth failed: %s", err))
	}
	decodedCreds := strings.Split(string(raw), ":")
	role, err := w.Handler.Login(decodedCreds[0], decodedCreds[1])
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Auth failed: %s", err))
	}
	token, err := CreateJWT(role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failure during JWT generation: %s", err))
	}
	return ctx.JSON(http.StatusOK, struct {
		Token string
	}{
		Token: token,
	})
}

func (w *ServerInterfaceWrapper) Signup(ctx echo.Context) error {
	decoder := json.NewDecoder(ctx.Request().Body)
	var user dto.User
	err := decoder.Decode(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
	}
	err = validator.Validate(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %s", err))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failure during password encryption: %s", err))
	}
	err = w.Handler.Signup(user.Username, string(hashedPassword))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failure during registration: %s", err))
	}
	return ctx.NoContent(http.StatusCreated)
}
