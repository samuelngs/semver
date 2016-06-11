package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/samuelngs/semver/backend"
	"github.com/samuelngs/semver/handler/v1"
)

// New creates server
func New(opts ...backend.Client) *gin.Engine {

	var store backend.Client

	for _, opt := range opts {
		store = opt
		break
	}

	if store == nil {
		log.Fatal("missing storage backend configuration")
	}

	api := gin.New()
	api.Use(gin.Recovery())

	// version 1
	v1.New(store, api)

	return api
}
