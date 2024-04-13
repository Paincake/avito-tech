package main

import (
	"context"
	"errors"
	"github.com/Paincake/avito-tech/internal/config"
	"github.com/Paincake/avito-tech/internal/database"
	"github.com/Paincake/avito-tech/internal/database/postgres"
	"github.com/Paincake/avito-tech/internal/server"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()

	configPath := os.Getenv("TEST_CONFIG_PATH")
	if configPath == "" {
		log.Fatal(errors.New("env variable not set"))
	}
	cfg, err := config.MustLoad(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db, _ := postgres.New(cfg.Name, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	done := make(chan bool)
	cache := server.NewMemoryCache(db, cfg.CacheKeyInvalidationTime, cfg.CacheSchedulerRate, done)
	defer close(done)
	ConfigureServer(db, cache, e, server.VerifyJWT, server.Logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down")
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func ConfigureServer(repository database.BannerRepository, cache server.BannerCache, e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	e.Use(middlewares...)
	si := server.Server{Repository: repository, Cache: cache}

	wrapper := server.ServerInterfaceWrapper{
		Handler: &si,
	}
	e.GET("/banner", wrapper.GetBanner)
	e.POST("/banner", wrapper.PostBanner)
	e.DELETE("/banner/:id", wrapper.DeleteBannerID)
	e.PATCH("/banner/:id", wrapper.PatchBannerID)
	e.GET("/user_banner", wrapper.GetUserBanner)
	e.POST("/login", wrapper.Login)
	e.POST("/signup", wrapper.Signup)
}
