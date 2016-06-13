package main

import (
	"github.com/samuelngs/semver/backend"
	"github.com/samuelngs/semver/pkg/env"
	"github.com/samuelngs/semver/server"
)

var defaultAddr = ":4000"

func main() {

	// create api server
	api := server.New(storage())

	// start server
	api.Run(defaultAddr)
}

func storage() backend.Client {
	switch env.Raw("SEMVER_BACKEND_STORAGE", "bolt") {
	case "bolt":
		return new(backend.Bolt)
	case "redis":
		return new(backend.Redis)
	case "cassandra":
		return new(backend.Cassandra)
	case "gce-datastore":
		return new(backend.GceDatastore)
	}
	return nil
}
