package orasclient

import (
	"context"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Interface is the oras interface.
type Interface interface {
	Push(ctx context.Context, resolver remotes.Resolver, ref string, provider content.Provider, descriptors []ocispec.Descriptor, opts ...oras.PushOpt) (ocispec.Descriptor, error)
	Pull(ctx context.Context, resolver remotes.Resolver, ref string, ingester content.Ingester, opts ...oras.PullOpt) (ocispec.Descriptor, []ocispec.Descriptor, error)
}
