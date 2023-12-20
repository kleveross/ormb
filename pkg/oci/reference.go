package oci

import (
	"fmt"

	"github.com/docker/distribution/reference"
)

type (
	// Reference defines the main components of a reference specification
	Reference struct {
		Tag  string
		Repo string
	}
)

// ParseReference converts a string to a Reference
func ParseReference(s string) (*Reference, error) {
	r, err := reference.Parse(s)
	if err != nil {
		return nil, err
	}

	var ref Reference

	if named, ok := r.(reference.Named); ok {
		ref.Repo = named.Name()
	}

	if tagged, ok := r.(reference.Tagged); ok {
		ref.Tag = tagged.Tag()
	}

	return &ref, nil
}

// FullName the full name of a reference (repo:tag)
func (ref *Reference) FullName() string {
	if ref.Tag == "" {
		return ref.Repo
	}
	return fmt.Sprintf("%s:%s", ref.Repo, ref.Tag)
}
