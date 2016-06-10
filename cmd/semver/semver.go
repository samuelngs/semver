package main

import (
	"github.com/samuelngs/semver/pkg/env"
	"github.com/samuelngs/semver/server"
)

var defaultAddr = ":4000"

func main() {

	// create api server
	api := server.New(env.Raw("SEMVER_BACKEND", "bolt"))

	// start server
	api.Listen(defaultAddr)
}
