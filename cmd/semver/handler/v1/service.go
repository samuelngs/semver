package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/redis.v3"

	"github.com/blang/semver"
	"github.com/labstack/echo"
	"github.com/samuelngs/semver/cmd/semver/backend"
	"github.com/satori/go.uuid"
)

const defaultVersion = "0.0.1"

func path(uuid string, dirs ...string) string {
	if uuid == "" {
		log.Fatal("uuid cannot be empty")
	}
	dir := fmt.Sprintf("semver:db:%s", uuid)
	for _, str := range dirs {
		dir += ":" + str
	}
	return dir
}

func errs(c echo.Context, e error) error {
	err := &Err{e.Error()}
	switch c.QueryParam("type") {
	case "xml":
		return c.XML(http.StatusForbidden, err)
	case "json":
		return c.JSON(http.StatusForbidden, err)
	default:
		return c.String(http.StatusForbidden, e.Error())
	}
}

func list(c echo.Context, i *List) error {
	switch c.QueryParam("type") {
	case "xml":
		return c.XML(http.StatusOK, i)
	case "json":
		return c.JSON(http.StatusOK, i)
	default:
		res := make([]string, len(i.Versions))
		for i, v := range i.Versions {
			res[i] = v.(string)
		}
		return c.String(http.StatusOK, strings.Join(res, "\n"))
	}
}

func resp(c echo.Context, i *Response) error {
	switch c.QueryParam("type") {
	case "xml":
		return c.XML(http.StatusOK, i)
	case "json":
		return c.JSON(http.StatusOK, i)
	default:
		if i.ID != "" {
			return c.String(http.StatusOK, i.ID)
		}
		return c.String(http.StatusOK, i.Version)
	}
}

// New Semver
func New(c echo.Context) error {
	var id string
	var version string
	if q := c.QueryParam("version"); q != "" {
		version = q
	} else {
		version = defaultVersion
	}
	v, err := semver.Make(version)
	if err != nil {
		return errs(c, errors.New("invalid semantic versioning format"))
	}
	for {
		id = uuid.NewV4().String()
		exists, err := backend.Exists(path(id))
		if err != nil {
			return err
		}
		if !exists {
			break
		}
	}
	if err := backend.Set(
		version,
		path(id, "version"),
		path(id, "archieve", version),
	); err != nil {
		return err
	}
	res := &Response{ID: id, Version: v.String(),
		Major: v.Major, Minor: v.Minor, Patch: v.Patch, Build: v.Build}
	return resp(c, res)
}

// Get Semver by `id`
func Get(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		return errs(c, errors.New("invalid semver uuid"))
	}
	s, err := backend.Get(path(id, "version"))
	if err != nil {
		if err == redis.Nil {
			return errs(c, errors.New("semver record does not exist in our database"))
		}
		return err
	}
	v, err := semver.Make(s)
	if err != nil {
		return errs(c, errors.New("invalid semantic versioning format"))
	}
	res := &Response{Version: v.String(), Major: v.Major, Minor: v.Minor, Patch: v.Patch, Build: v.Build}
	return resp(c, res)
}

// Set Semver by `id`
func Set(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		return errs(c, errors.New("invalid semver uuid"))
	}
	v, err := semver.Make(c.FormValue("version"))
	if err != nil {
		return errs(c, errors.New("invalid semantic versioning format"))
	}
	exists, err := backend.Exists(path(id, "version"))
	if err != nil {
		return err
	} else if !exists {
		return errs(c, errors.New("semver record does not exist in our database"))
	}
	if err := backend.Set(
		v.String(),
		path(id, "version"),
		path(id, "archieve", v.String()),
	); err != nil {
		return err
	}
	res := &Response{Version: v.String(), Major: v.Major, Minor: v.Minor, Patch: v.Patch, Build: v.Build}
	return resp(c, res)
}

// Bump version by type {major, minor, patch}
func Bump(c echo.Context) (e error) {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		return errs(c, errors.New("invalid semver uuid"))
	}
	s, err := backend.Get(path(id, "version"))
	if err != nil {
		if err == redis.Nil {
			return errs(c, errors.New("semver record does not exist in our database"))
		}
		return err
	}
	v, err := semver.Make(s)
	if err != nil {
		return errs(c, errors.New("invalid semantic versioning format"))
	}
	switch c.QueryParam("level") {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case "minor":
		v.Minor++
		v.Patch = 0
	default:
		v.Patch++
	}
	v.Pre = make([]semver.PRVersion, 0)
	if err := v.Validate(); err != nil {
		return errs(c, errors.New("invalid semantic versioning format"))
	}
	if err := backend.Set(
		v.String(),
		path(id, "version"),
		path(id, "archieve", v.String()),
	); err != nil {
		return err
	}
	res := &Response{Version: v.String(), Major: v.Major, Minor: v.Minor, Patch: v.Patch, Build: v.Build}
	return resp(c, res)
}

// Delete semver data
func Delete(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		return errs(c, errors.New("invalid semver uuid"))
	}
	exists, err := backend.Exists(path(id, "version"))
	if err != nil {
		return err
	} else if !exists {
		return errs(c, errors.New("semver record does not exist in our database"))
	}
	if err := backend.Delete(path(id, "*")); err != nil {
		return err
	}
	res := &Response{ID: id}
	return resp(c, res)
}

// History to list semver records
func History(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.FromString(id); err != nil {
		return errs(c, errors.New("invalid semver uuid"))
	}
	exists, err := backend.Exists(path(id, "version"))
	if err != nil {
		return err
	} else if !exists {
		return errs(c, errors.New("semver record does not exist in our database"))
	}
	refs, err := backend.List(path(id, "archieve:*"))
	if err != nil {
		return errs(c, errors.New("semver record does not exist in our database"))
	}
	vals, err := backend.GetAll(refs...)
	if err != nil {
		return errs(c, errors.New("semver record does not exist in our database"))
	}
	res := &List{vals}
	return list(c, res)
}
