package server

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/samuelngs/semver/backend"
	"github.com/samuelngs/semver/handler/v1"
)

// New creates server
func New(opts ...string) *iris.Iris {

	var s string
	for _, opt := range opts {
		s = opt
		break
	}

	conf := config.Iris{
		Profile:               false,
		DisablePathCorrection: true,
		DisableBanner:         true,
	}

	api := iris.New(conf)

	// create backend
	store := backend.Get(s)

	// version 1
	v1.New(store, api)

	return api
}
