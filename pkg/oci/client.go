package oci

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/caicloud/ormb/pkg/consts"
	"github.com/caicloud/ormb/pkg/model"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

const (
	credentialsFileBasename = "config.json"
)

// Client works with OCI-compliant registries and local cache.
type Client struct {
	debug      bool
	out        io.Writer
	authorizer *Authorizer
	resolver   *Resolver
	cache      *Cache
	rootPath   string
	plainHTTP  bool
}

// NewClient returns a new registry client with config
func NewClient(opts ...ClientOption) (*Client, error) {
	client := &Client{
		out: ioutil.Discard,
	}
	for _, opt := range opts {
		opt(client)
	}
	// set defaults if fields are missing
	if client.authorizer == nil {
		credentialsFile := path.Join(client.rootPath, credentialsFileBasename)
		authClient, err := auth.NewClient(credentialsFile)
		if err != nil {
			return nil, err
		}
		client.authorizer = &Authorizer{
			Client: authClient,
		}
	}
	if client.resolver == nil {
		resolver, err := client.authorizer.Resolver(context.Background(), http.DefaultClient, client.plainHTTP)
		if err != nil {
			return nil, err
		}
		client.resolver = &Resolver{
			Resolver: resolver,
		}
	}
	if client.cache == nil {
		cache, err := NewCache(
			CacheOptDebug(client.debug),
			CacheOptWriter(client.out),
			CacheOptRoot(path.Join(client.rootPath, cacheRootDir)),
		)
		if err != nil {
			return nil, err
		}
		client.cache = cache
	}
	return client, nil
}

// Login logs into a registry
func (c *Client) Login(hostname string, username string, password string, insecure bool) error {
	err := c.authorizer.Login(ctx(c.out, c.debug), hostname, username, password, insecure)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Login succeeded\n")
	return nil
}

// Logout logs out of a registry
func (c *Client) Logout(hostname string) error {
	err := c.authorizer.Logout(ctx(c.out, c.debug), hostname)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.out, "Logout succeeded")
	return nil
}

// SaveModel stores a copy of model in local cache
func (c *Client) SaveModel(ch *model.Model, ref *Reference) error {
	r, err := c.cache.StoreReference(ref, ch)
	if err != nil {
		return err
	}
	c.printCacheRefSummary(r)

	// Store the manifest in index.json in local cache.
	err = c.cache.AddManifest(ref, r.Manifest)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%s: saved\n", r.Tag)
	return nil
}

// PushModel uploads a model to a registry.
func (c *Client) PushModel(ref *Reference) error {
	r, err := c.cache.FetchReference(ref)
	if err != nil {
		return err
	}
	if !r.Exists {
		return errors.New(fmt.Sprintf("Model not found: %s", r.Name))
	}
	fmt.Fprintf(c.out, "The push refers to repository [%s]\n", r.Repo)
	c.printCacheRefSummary(r)
	layers := []ocispec.Descriptor{*r.ContentLayer}
	_, err = oras.Push(ctx(c.out, c.debug), c.resolver, r.Name, c.cache.Provider(), layers,
		oras.WithConfig(*r.Config), oras.WithNameValidation(nil))
	if err != nil {
		return err
	}
	s := ""
	numLayers := len(layers)
	if 1 < numLayers {
		s = "s"
	}
	fmt.Fprintf(c.out,
		"%s: pushed to remote (%d layer%s, %s total)\n", r.Tag, numLayers, s, byteCountBinary(r.Size))
	return nil
}

// RemoveModel deletes a locally saved model.
func (c *Client) RemoveModel(ref *Reference) error {
	r, err := c.cache.DeleteReference(ref)
	if err != nil {
		return err
	}
	if !r.Exists {
		return errors.New(fmt.Sprintf("Model not found: %s", ref.FullName()))
	}
	fmt.Fprintf(c.out, "%s: removed\n", r.Tag)
	return nil
}

// PullModel downloads a model from a registry.
func (c *Client) PullModel(ref *Reference) error {
	if ref.Tag == "" {
		return errors.New("tag explicitly required")
	}
	existing, err := c.cache.FetchReference(ref)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%s: Pulling from %s\n", ref.Tag, ref.Repo)
	manifest, _, err := oras.Pull(ctx(c.out, c.debug), c.resolver, ref.FullName(), c.cache.Ingester(),
		oras.WithPullEmptyNameAllowed(),
		oras.WithAllowedMediaTypes(consts.KnownMediaTypes()),
		oras.WithContentProvideIngester(c.cache.ProvideIngester()))
	if err != nil {
		return err
	}
	err = c.cache.AddManifest(ref, &manifest)
	if err != nil {
		return err
	}
	r, err := c.cache.FetchReference(ref)
	if err != nil {
		return err
	}
	if !r.Exists {
		return errors.New(fmt.Sprintf("Model not found: %s", r.Name))
	}
	c.printCacheRefSummary(r)
	if !existing.Exists {
		fmt.Fprintf(c.out, "Status: Downloaded newer model for %s\n", ref.FullName())
	} else {
		fmt.Fprintf(c.out, "Status: Model is up to date for %s\n", ref.FullName())
	}
	return err
}

// LoadModel retrieves a model object by reference
func (c *Client) LoadModel(ref *Reference) (*model.Model, error) {
	r, err := c.cache.FetchReference(ref)
	if err != nil {
		return nil, err
	}
	if !r.Exists {
		return nil, errors.New(fmt.Sprintf("Model not found: %s", ref.FullName()))
	}
	c.printCacheRefSummary(r)
	return r.Model, nil
}

// printCacheRefSummary prints out model ref summary
func (c *Client) printCacheRefSummary(r *CacheRefSummary) {
	fmt.Fprintf(c.out, "ref:     %s\n", r.Name)
	fmt.Fprintf(c.out, "digest:  %s\n", r.Digest.Hex())
	fmt.Fprintf(c.out, "size:    %s\n", byteCountBinary(r.Size))
	// fmt.Fprintf(c.out, "name:    %s\n", r.Model.Metadata.Name)
	// fmt.Fprintf(c.out, "version: %s\n", r.Model.Metadata.Version)
}
