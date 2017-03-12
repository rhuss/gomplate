package vault

import (
	"io/ioutil"
	"os"

	"github.com/blang/vfs"
)

// GetValue - gets a value from either an environment variable, or if it isn't
// set, from a path specified by another environment variable named the same as
// the original variable, but suffixed with `_FILE`.
// If a `*_FILE` environment variable is set, but the file doesn't exist, this
// function panics.
func GetValue(key string, fs vfs.Filesystem) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	p := os.Getenv(key + "_FILE")
	if p == "" {
		return ""
	}
	f, err := fs.OpenFile(p, os.O_RDONLY, 0)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(b)
}
