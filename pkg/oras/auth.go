package oras

import (
	"github.com/deislabs/oras/pkg/auth"
)

type (
	// Authorizer handles registry auth operations
	Authorizer struct {
		auth.Client
	}
)
