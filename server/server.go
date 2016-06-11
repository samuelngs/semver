package server

import (
	"github.com/gin-gonic/gin"
	"github.com/samuelngs/semver/backend"
	"github.com/samuelngs/semver/handler/v1"
	"github.com/samuelngs/semver/pkg/env"
)

// New creates server
func New(opts ...string) *gin.Engine {

	s := env.Raw("SEMVER_BACKEND_STORAGE", "bolt")
	for _, opt := range opts {
		s = opt
		break
	}

	api := gin.New()
	api.Use(gin.Recovery())

	// create backend
	store := backend.Get(s)

	// version 1
	v1.New(store, api)

	return api
}
