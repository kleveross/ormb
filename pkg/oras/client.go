package oras

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/deislabs/oras/pkg/oras"
	"github.com/kleveross/ormb/pkg/consts"
	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/oci"
	"github.com/kleveross/ormb/pkg/oras/cache"
	"github.com/kleveross/ormb/pkg/oras/orasclient"
	bts "github.com/kleveross/ormb/pkg/util/bytes"
	"github.com/kleveross/ormb/pkg/util/ctx"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/pkg/errors"
	"github.com/xeonx/timeago"
)

const (
	credentialsFileBasename = "config.json"
	dockerCredentialsDirectory = "~/.docker/"
)

var _ Interface = (*Client)(nil)

var credentialsFile string

// Client works with OCI-compliant registries and local cache.
type Client struct {
	debug      bool
	out        io.Writer
	authorizer *Authorizer
	resolver   *Resolver
	cache      cache.Interface
	orasClient orasclient.Interface
	rootPath   string
	plainHTTP  bool
	insecure   bool
}

// NewClient returns a new registry client with config
func NewClient(opts ...ClientOption) (Interface, error) {
	client := &Client{
		out:        ioutil.Discard,
		orasClient: orasclient.New(),
	}
	for _, opt := range opts {
		opt(client)
	}
	// set defaults if fields are missing
	if client.authorizer == nil {
	
		if fileExists(path.Join((dockerCredentialsDirectory))) {
			credentialsFile = path.Join(dockerCredentialsDirectory, credentialsFileBasename)	
			fmt.Println("Using Docker Config for Login")
		} else {
			credentialsFile = path.Join(client.rootPath, credentialsFileBasename)
			fmt.Println("Using ORMB Config for Login")
		}
		authClient, err := auth.NewClient(credentialsFile)
		if err != nil {
			return nil, err
		}
		client.authorizer = &Authorizer{
			Client: authClient,
		}
	}
	if client.resolver == nil {
		resolver, err := client.authorizer.Resolver(
			context.Background(),
			&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: client.insecure},
				},
			}, client.plainHTTP)
		if err != nil {
			return nil, err
		}
		client.resolver = &Resolver{
			Resolver: resolver,
		}
	}
	if client.cache == nil {
		cache, err := cache.New(
			cache.CacheOptDebug(client.debug),
			cache.CacheOptWriter(client.out),
			cache.CacheOptRoot(path.Join(client.rootPath, cache.CacheRootDir)),
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
	if insecure {
		fmt.Fprintf(c.out, "Login insecurely\n")
	}
	err := c.authorizer.Login(ctx.Context(c.out, c.debug), hostname, username, password, insecure)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "Login succeeded\n")
	return nil
}

// Logout logs out of a registry
func (c *Client) Logout(hostname string) error {
	err := c.authorizer.Logout(ctx.Context(c.out, c.debug), hostname)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.out, "Logout succeeded")
	return nil
}

// SaveModel stores a copy of model in local cache
func (c *Client) SaveModel(ch *model.Model, ref *oci.Reference) error {
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

// TagModel tags the ref to target.
func (c *Client) TagModel(ref *oci.Reference, target *oci.Reference) error {
	if ref.FullName() == target.FullName() {
		return nil
	}
	if err := c.cache.TagReference(ref, target); err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%s: tagged\n", ref.FullName())
	return nil
}

// PushModel uploads a model to a registry.
func (c *Client) PushModel(ref *oci.Reference) error {
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
	_, err = c.orasClient.Push(ctx.Context(c.out, c.debug),
		c.resolver, r.Name, c.cache.Provider(), layers,
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
		"%s: pushed to remote (%d layer%s, %s total)\n", r.Tag, numLayers, s,
		bts.ByteCountBinary(r.Size))
	return nil
}

// RemoveModel deletes a locally saved model.
func (c *Client) RemoveModel(ref *oci.Reference) error {
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
func (c *Client) PullModel(ref *oci.Reference) error {
	if ref.Tag == "" {
		return errors.New("tag explicitly required")
	}
	existing, err := c.cache.FetchReference(ref)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.out, "%s: Pulling from %s\n", ref.Tag, ref.Repo)
	manifest, _, err := c.orasClient.Pull(ctx.Context(c.out, c.debug),
		c.resolver, ref.FullName(), c.cache.Ingester(),
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
func (c *Client) LoadModel(ref *oci.Reference) (*model.Model, error) {
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

func (c *Client) Models() error {
	refs, err := c.cache.ListReferences()
	if err != nil {
		return err
	}

	maxLen := []int{0, 0, 12, 20, 0}
	for _, ref := range refs {
		if len(ref.Repo) > maxLen[0] {
			maxLen[0] = len(ref.Repo)
		}
		if len(ref.Tag) > maxLen[1] {
			maxLen[1] = len(ref.Tag)
		}

		size := bts.ByteCountBinary(ref.Size)
		if len(size) > maxLen[4] {
			maxLen[4] = len(size)
		}
	}

	format := "%-" + strconv.Itoa(maxLen[0]) + "s "  // repo name
	format += "%-" + strconv.Itoa(maxLen[1]) + "s "  // tag
	format += "%-" + strconv.Itoa(maxLen[2]) + "s "  // digest
	format += "%-" + strconv.Itoa(maxLen[3]) + "s "  // create at
	format += "%-" + strconv.Itoa(maxLen[4]) + "s\n" // size
	fmt.Fprintf(c.out, format, "REPOSITORY", "TAG", "MODEL ID", "CREATED", "SIZE")

	for _, ref := range refs {
		fmt.Fprintf(c.out, format, ref.Repo, ref.Tag, ref.Digest.Hex()[0:12], timeago.English.Format(ref.CreatedAt), bts.ByteCountBinary(ref.Size))
	}
	return nil
}

// printCacheRefSummary prints out model ref summary
func (c *Client) printCacheRefSummary(r *cache.CacheRefSummary) {
	fmt.Fprintf(c.out, "ref:       %s\n", r.Name)
	fmt.Fprintf(c.out, "digest:    %s\n", r.Digest.Hex())
	fmt.Fprintf(c.out, "size:      %s\n", bts.ByteCountBinary(r.Size))
	if r.Model != nil && r.Model.Metadata != nil {
		fmt.Fprintf(c.out, "format:    %s\n", r.Model.Metadata.Format)
	}
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
