package oci

import (
	"io"
)

type (
	// CacheOption allows specifying various settings configurable by the user for overriding the defaults
	// used when creating a new default cache
	CacheOption func(*Cache)
)

// CacheOptDebug returns a function that sets the debug setting on cache options set
func CacheOptDebug(debug bool) CacheOption {
	return func(cache *Cache) {
		cache.debug = debug
	}
}

// CacheOptWriter returns a function that sets the writer setting on cache options set
func CacheOptWriter(out io.Writer) CacheOption {
	return func(cache *Cache) {
		cache.out = out
	}
}

// CacheOptRoot returns a function that sets the root directory setting on cache options set
func CacheOptRoot(rootDir string) CacheOption {
	return func(cache *Cache) {
		cache.rootDir = rootDir
	}
}
