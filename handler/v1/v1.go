package v1

import (
	"github.com/kataras/iris"
	"github.com/samuelngs/semver/backend"
)

const defaultVersion = "0.0.1"

// New create route
func New(b backend.Client, c *iris.Iris) *Router {

	r := &Router{backend.New(b)}

	// GET: /v1/new
	c.Get("/new", r.Create)

	g := c.Party("/v1")
	{
		// GET: /v1
		g.Get("", r.Default)
		g.Get("/", r.Default)

		// GET: /v1/{project-id}/history
		g.Get("/:id/history", r.History)

		// GET: /v1/{project-id}/bump
		g.Get("/:id/bump", r.Bump)

		// GET: /v1/{project-id}
		g.Get("/:id", r.Get)

		// POST: /v1/{project-id}
		g.Post("/:id", r.Set)

		// DELETE: /v1/{project-id}
		g.Delete("/:id", r.Delete)
	}
	return r
}
