package v1

import (
	"net/http"
	"strings"

	"github.com/blang/semver"
	"github.com/gin-gonic/gin"
	"github.com/samuelngs/semver/backend"
	"github.com/satori/go.uuid"
)

// Router route
type Router struct {
	m *backend.Manager
}

// Uniq generate unique id
func (r *Router) uniq() (string, error) {
	var id string
	for {
		id = uuid.NewV4().String()
		exists, err := r.m.Exists(r.m.Path(id))
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
	}
	return id, nil
}

// version semver string
func (r *Router) version(s string) (semver.Version, error) {
	v, err := semver.Make(s)
	if err != nil {
		return v, ErrInvalidVersioningFormat
	}
	return v, nil
}

// uuid parser
func (r *Router) uuid(s string) (uuid.UUID, error) {
	uuid, err := uuid.FromString(s)
	if err != nil {
		return uuid, ErrInvalidUUID
	}
	return uuid, nil
}

func (r *Router) release(c *gin.Context) {
	if o := recover(); o != nil {
		r.err(c, ErrInternalServer)
	}
}

// Err prints error message
func (r *Router) err(c *gin.Context, e error) {
	switch c.DefaultQuery("output", "text") {
	case "xml":
		c.XML(http.StatusForbidden, &Warning{e.Error()})
	case "json":
		c.JSON(http.StatusForbidden, &Warning{e.Error()})
	default:
		c.String(http.StatusForbidden, "%v", e)
	}
}

// Echo prints data message
func (r *Router) echo(c *gin.Context, d interface{}) {
	switch c.DefaultQuery("output", "text") {
	case "xml":
		c.XML(http.StatusOK, d)
	case "json":
		c.JSON(http.StatusOK, d)
	default:
		c.String(http.StatusOK, "%v", d)
	}
}

// Default route
func (r *Router) Default(c *gin.Context) {
	defer r.release(c)
	c.String(http.StatusOK, "ok")
}

// Create is the new semver handler
func (r *Router) Create(c *gin.Context) {
	defer r.release(c)
	var s string
	if v := strings.TrimSpace(c.Query("version")); v != "" {
		s = v
	} else {
		s = defaultVersion
	}
	ver, err := r.version(s)
	if err != nil {
		r.err(c, err)
		return
	}
	id, err := r.uniq()
	if err != nil {
		r.err(c, err)
		return
	}
	if err := r.m.Set(
		ver.String(),
		r.m.Path(id, "version"),
		r.m.Path(id, "archive", ver.String()),
	); err != nil {
		r.err(c, err)
		return
	}
	res := &Versioning{
		Project: id,
		Version: ver.String(),
		Major:   ver.Major,
		Minor:   ver.Minor,
		Patch:   ver.Patch,
		Build:   ver.Build,
	}
	r.echo(c, res)
}

// Get semver by project `id`
func (r *Router) Get(c *gin.Context) {
	defer r.release(c)
	id := c.Param("id")
	if id == "new" {
		r.Create(c)
		return
	}
	if _, err := r.uuid(id); err != nil {
		r.err(c, err)
		return
	}
	vers, err := r.m.Get(
		r.m.Path(id, "version"),
	)
	if err != nil || len(vers) <= 0 {
		r.err(c, ErrProjectNotFound)
		return
	}
	ver, err := r.version(vers[0])
	if err != nil {
		r.err(c, err)
		return
	}
	res := &Versioning{
		Version: ver.String(),
		Major:   ver.Major,
		Minor:   ver.Minor,
		Patch:   ver.Patch,
		Build:   ver.Build,
	}
	r.echo(c, res)
}

// Set Semver by `id`
func (r *Router) Set(c *gin.Context) {
	defer r.release(c)
	id := c.Param("id")
	if _, err := r.uuid(id); err != nil {
		r.err(c, err)
		return
	}
	exists, err := r.m.Exists(r.m.Path(id, "version"))
	if err != nil {
		r.err(c, err)
		return
	} else if !exists {
		r.err(c, ErrProjectNotFound)
		return
	}
	ver, err := r.version(c.PostForm("version"))
	if err != nil {
		r.err(c, err)
		return
	}
	if err := r.m.Set(
		ver.String(),
		r.m.Path(id, "version"),
		r.m.Path(id, "archive", ver.String()),
	); err != nil {
		r.err(c, err)
		return
	}
	res := &Versioning{
		Version: ver.String(),
		Major:   ver.Major,
		Minor:   ver.Minor,
		Patch:   ver.Patch,
		Build:   ver.Build,
	}
	r.echo(c, res)
}

// Bump version by type {major, minor, patch}
func (r *Router) Bump(c *gin.Context) {
	defer r.release(c)
	id := c.Param("id")
	if _, err := r.uuid(id); err != nil {
		r.err(c, err)
		return
	}
	exists, err := r.m.Exists(r.m.Path(id, "version"))
	if err != nil {
		r.err(c, err)
		return
	} else if !exists {
		r.err(c, ErrProjectNotFound)
		return
	}
	vers, err := r.m.Get(
		r.m.Path(id, "version"),
	)
	if len(vers) <= 0 {
		r.err(c, ErrProjectNotFound)
		return
	}
	ver, err := r.version(vers[0])
	if err != nil {
		r.err(c, err)
		return
	}
	if typ := c.Query("type"); typ == "major" {
		ver.Major++
		ver.Minor = 0
		ver.Patch = 0
	} else if typ == "minor" {
		ver.Minor++
		ver.Patch = 0
	} else {
		ver.Patch++
	}
	ver.Pre = make([]semver.PRVersion, 0)
	if err := ver.Validate(); err != nil {
		r.err(c, err)
		return
	}
	if err := r.m.Set(
		ver.String(),
		r.m.Path(id, "version"),
		r.m.Path(id, "archive", ver.String()),
	); err != nil {
		r.err(c, err)
		return
	}
	res := &Versioning{
		Version: ver.String(),
		Major:   ver.Major,
		Minor:   ver.Minor,
		Patch:   ver.Patch,
		Build:   ver.Build,
	}
	r.echo(c, res)
}

// History to list semver records
func (r *Router) History(c *gin.Context) {
	defer r.release(c)
	id := c.Param("id")
	if _, err := r.uuid(id); err != nil {
		r.err(c, err)
		return
	}
	exists, err := r.m.Exists(r.m.Path(id, "version"))
	if err != nil {
		r.err(c, err)
		return
	} else if !exists {
		r.err(c, ErrProjectNotFound)
		return
	}
	keys, err := r.m.List(r.m.Path(id, "archive"))
	if err != nil {
		r.err(c, err)
		return
	}
	vers, err := r.m.Get(keys...)
	if err != nil {
		r.err(c, err)
		return
	}
	arch := &Archive{
		Versions: make([]*Versioning, len(vers)),
	}
	for i, s := range vers {
		ver, err := r.version(s)
		if err != nil {
			r.err(c, err)
			return
		}
		arch.Versions[i] = &Versioning{
			Version: ver.String(),
			Major:   ver.Major,
			Minor:   ver.Minor,
			Patch:   ver.Patch,
			Build:   ver.Build,
		}
	}
	r.echo(c, arch)
}

// Delete to remove project
func (r *Router) Delete(c *gin.Context) {
	defer r.release(c)
	id := c.Param("id")
	if _, err := r.uuid(id); err != nil {
		r.err(c, err)
		return
	}
	exists, err := r.m.Exists(r.m.Path(id, "version"))
	if err != nil {
		r.err(c, err)
		return
	} else if !exists {
		r.err(c, ErrProjectNotFound)
		return
	}
	keys, err := r.m.List(r.m.Path(id))
	if err != nil {
		r.err(c, err)
		return
	}
	if err := r.m.Delete(keys...); err != nil {
		r.err(c, err)
		return
	}
	c.String(http.StatusOK, "ok")
}
