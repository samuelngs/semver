package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/samuelngs/semver/backend"
)

const defaultVersion = "0.0.1"

// New create route
func New(b backend.Client, c *gin.Engine) *Router {

	r := &Router{backend.New(b)}

	g := c.Group("/v1")
	{
		// GET: /v1
		g.GET("", r.Default)
		g.GET("/", r.Default)

		// GET: /v1/{project-id} or /v1/new
		g.GET("/:id", r.Get)

		// POST: /v1/{project-id}
		g.POST("/:id", r.Set)

		// DELETE: /v1/{project-id}
		g.DELETE("/:id", r.Delete)

		// GET: /v1/{project-id}/history
		g.GET("/:id/history", r.History)

		// GET: /v1/{project-id}/bump
		g.GET("/:id/bump", r.Bump)
	}
	return r
}
