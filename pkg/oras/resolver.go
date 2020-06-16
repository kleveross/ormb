package oras

import (
	"github.com/containerd/containerd/remotes"
)

type (
	// Resolver provides remotes based on a locator
	Resolver struct {
		remotes.Resolver
	}
)
