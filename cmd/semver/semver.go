package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/samuelngs/semver/cmd/semver/handler/v1"
)

var defaultAddr = ":4000"

func main() {

	e := echo.New()

	e.Get("/", v1.Default)
	e.Get("/v1/new", v1.New)
	e.Get("/v1/:id", v1.Get)
	e.Get("/v1/:id/bump", v1.Bump)
	e.Get("/v1/:id/history", v1.History)
	e.Post("/v1/:id", v1.Set)
	e.Delete("/v1/:id", v1.Delete)

	e.Run(standard.New(defaultAddr))
}
