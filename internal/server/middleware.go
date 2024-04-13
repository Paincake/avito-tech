package server

import (
	"context"
	"fmt"
	"github.com/Paincake/avito-tech/internal/dto"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	f := middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		HandleError:  true,
		LogLatency:   true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	})
	return f(next)
}

func VerifyJWT(next echo.HandlerFunc) echo.HandlerFunc {
	var role string
	return func(c echo.Context) error {
		if c.Request().RequestURI == "/login" || c.Request().RequestURI == "/signup" {
			return next(c)
		}
		r := c.Request().Header
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		if r.Get("Token") != "" {
			token, err := jwt.Parse(r.Get("Token"), func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil {
				logger.Debug("Request discarded: auth failed")
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Request discarded: auth failed: %s", err))
			}
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
					role = claims["role"].(string)
					c.Set(dto.TokenRoleContextKey, role)
					logger.Debug(fmt.Sprintf("User with claims %s authenticated", role))
				} else {
					logger.Debug("Request discarded: auth failed: required claim absent")
					return echo.NewHTTPError(http.StatusUnauthorized, "Request discarded: auth failed: required claim absent")
				}
			} else {
				logger.Debug("Request discarded: auth failed: invalid token")
				return echo.NewHTTPError(http.StatusUnauthorized, "Request discarded: auth failed: invalid token")
			}
		} else {
			logger.Debug("Request discarded: auth failed: token absent")
			return echo.NewHTTPError(http.StatusBadRequest, "request discarded: Token header parameter absent")
		}
		return next(c)
	}
}
