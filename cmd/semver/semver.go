package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/samuelngs/semver/backend"
	"github.com/samuelngs/semver/handler/v1"
)

var defaultAddr = ":4000"

func main() {

	conf := config.Iris{
		Profile:               false,
		DisablePathCorrection: true,
		DisableBanner:         true,
	}

	api := iris.New(conf)

	// version 1
	v1.New(new(backend.Bolt), api)

	// start api server
	api.Listen(defaultAddr)
}
