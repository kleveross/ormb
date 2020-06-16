package orasclient

import (
	"context"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

var _ Interface = (*Client)(nil)

// Client is the client for oras.
type Client struct {
}

func New() Interface {
	return &Client{}
}

// Push pushes files to the remote.
func (o Client) Push(ctx context.Context, resolver remotes.Resolver, ref string, provider content.Provider, descriptors []ocispec.Descriptor, opts ...oras.PushOpt) (ocispec.Descriptor, error) {
	return oras.Push(ctx, resolver, ref, provider, descriptors, opts...)
}

// Pull pulls files from the remote.
func (o Client) Pull(ctx context.Context, resolver remotes.Resolver, ref string, ingester content.Ingester, opts ...oras.PullOpt) (ocispec.Descriptor, []ocispec.Descriptor, error) {
	return oras.Pull(ctx, resolver, ref, ingester, opts...)
}
