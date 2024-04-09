package main

import (
	"github.com/Paincake/avito-tech/internal/server"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Use(server.VerifyJWT)
	e.Use(server.Logger)

}
