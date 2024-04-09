package server

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
)

const TokenRoleContextKey = "Token"

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			HandleError: true,
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
		return next(c)
	}
}

func VerifyJWT(next echo.HandlerFunc) echo.HandlerFunc {
	var role string
	return func(c echo.Context) error {
		logger := c.Logger()
		r := c.Request().Header
		if r.Get("Token") != "" {
			token, err := jwt.Parse(r.Get("Token"), func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil {
				c.Response().WriteHeader(http.StatusUnauthorized)
				c.Error(fmt.Errorf("invalid token: %s", err))
				logger.Debug("Request discarded: auth failed")
				return fmt.Errorf("invalid token: %s", err)
			}
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
					role = claims["role"].(string)
					c.Set(TokenRoleContextKey, role)
					logger.Debug(fmt.Sprintf("User with claims %s authenticated", role))
				} else {
					logger.Debug("Request discarded: auth failed: required claim absent")
					c.Response().WriteHeader(http.StatusUnauthorized)
					c.Error(fmt.Errorf("request discarded: auth failed: required claim absent"))
					return fmt.Errorf("request discarded: auth failed: required claim absent")
				}
			} else {
				c.Response().WriteHeader(http.StatusUnauthorized)
				logger.Debug("Request discarded: auth failed: invalid token")
				return fmt.Errorf("request discarded: auth failed: invalid token")
			}
		} else {
			logger.Debug("Request discarded: auth failed: token absent")
			c.Error(fmt.Errorf("request discarded: auth failed: token absent"))
			return fmt.Errorf("request discarded: auth failed: token absent")
		}
		return next(c)
	}
}
