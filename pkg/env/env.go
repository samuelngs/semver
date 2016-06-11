package env

import (
	"os"
	"strconv"
)

// Set env
func Set(key string, vals ...string) error {
	var val string
	for _, v := range vals {
		val = v
		break
	}
	return os.Setenv(key, val)
}

// Raw to read environment key
func Raw(name string, defs ...string) string {
	var def string
	for _, d := range defs {
		def = d
		break
	}
	if str := os.Getenv(name); str != "" {
		return str
	}
	return def
}

// I64 to read environment key and return value in int64 format
func I64(name string, defs ...int64) int64 {
	var def int64
	for _, d := range defs {
		def = d
		break
	}
	v := Raw(name)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return i
}

// Int to read environment key and return value in int format
func Int(name string, defs ...int) int {
	var def int
	for _, d := range defs {
		def = d
		break
	}
	v := Raw(name)
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
