package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	orascontent "github.com/deislabs/oras/pkg/content"
	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"

	"github.com/caicloud/ormb/pkg/consts"
	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
	"github.com/caicloud/ormb/pkg/parser"
	"github.com/caicloud/ormb/pkg/util/ctx"
)

const (
	// CacheRootDir is the root directory for a cache
	CacheRootDir = "cache"
)

var _ Interface = (*Cache)(nil)

// Cache handles local/in-memory storage of Helm charts, compliant with OCI Layout
type Cache struct {
	debug   bool
	out     io.Writer
	rootDir string
	parser  parser.Parser

	ociStore *orascontent.OCIStore
	// TODO(gaocegege): Do we really need it?
	memoryStore *orascontent.Memorystore
}

// CacheRefSummary contains as much info as available describing a chart reference in cache
// Note: fields here are sorted by the order in which they are set in FetchReference method
type CacheRefSummary struct {
	Name         string
	Repo         string
	Tag          string
	Exists       bool
	Manifest     *ocispec.Descriptor
	Config       *ocispec.Descriptor
	ContentLayer *ocispec.Descriptor
	Size         int64
	Digest       digest.Digest
	CreatedAt    time.Time
	Model        *model.Model
}

// New returns a new OCI Layout-compliant cache with config
func New(opts ...CacheOption) (Interface, error) {
	cache := &Cache{
		out:    ioutil.Discard,
		parser: parser.NewDefaultParser(),
	}
	for _, opt := range opts {
		opt(cache)
	}
	// validate
	if cache.rootDir == "" {
		return nil, errors.New("must set cache root dir on initialization")
	}
	return cache, nil
}

// FetchReference retrieves a model ref from cache.
func (cache *Cache) FetchReference(ref *oci.Reference) (*CacheRefSummary, error) {
	if err := cache.init(); err != nil {
		return nil, err
	}
	r := CacheRefSummary{
		Name: ref.FullName(),
		Repo: ref.Repo,
		Tag:  ref.Tag,
	}
	for _, desc := range cache.ociStore.ListReferences() {
		if desc.Annotations[ocispec.AnnotationRefName] == r.Name {
			r.Exists = true

			// Fetch the manifest.
			manifestBytes, err := cache.fetchBlob(&desc)
			if err != nil {
				return &r, err
			}
			var manifest ocispec.Manifest
			err = json.Unmarshal(manifestBytes, &manifest)
			if err != nil {
				return &r, err
			}
			r.Manifest = &desc

			// TODO(gaocegege): Fetch the config.
			r.Config = &manifest.Config
			numLayers := len(manifest.Layers)
			if numLayers != 1 {
				return &r, errors.New(
					fmt.Sprintf("manifest does not contain exactly 1 layer (total: %d)", numLayers))
			}

			// Fetch the content.
			var contentLayer *ocispec.Descriptor
			for _, layer := range manifest.Layers {
				switch layer.MediaType {
				case consts.MediaTypeModelContentLayer:
					contentLayer = &layer
				}
			}
			if contentLayer == nil {
				return &r, errors.New(
					fmt.Sprintf(
						"manifest does not contain a layer with mediatype %s",
						consts.MediaTypeModelContentLayer))
			}
			if contentLayer.Size == 0 {
				return &r, errors.New(
					fmt.Sprintf(
						"manifest layer with mediatype %s is of size 0",
						consts.MediaTypeModelContentLayer))
			}
			r.ContentLayer = contentLayer
			info, err := cache.ociStore.Info(ctx.Context(cache.out, cache.debug), contentLayer.Digest)
			if err != nil {
				return &r, err
			}
			r.Size = info.Size
			r.Digest = info.Digest
			r.CreatedAt = info.CreatedAt
			contentBytes, err := cache.fetchBlob(contentLayer)
			if err != nil {
				return &r, err
			}
			configBytes, err := cache.fetchBlob(r.Config)
			if err != nil {
				return &r, err
			}
			metadata, err := cache.parser.Parse(configBytes)
			if err != nil {
				return &r, err
			}

			// TODO(gaocegege): Optimize the memory usage.
			r.Model = &model.Model{
				Content:  contentBytes,
				Config:   configBytes,
				Metadata: metadata,
			}
		}
	}
	return &r, nil
}

// StoreReference stores a model ref in cache
func (cache *Cache) StoreReference(ref *oci.Reference, m *model.Model) (*CacheRefSummary, error) {
	if err := cache.init(); err != nil {
		return nil, err
	}
	r := CacheRefSummary{
		Name:  ref.FullName(),
		Repo:  ref.Repo,
		Tag:   ref.Tag,
		Model: m,
	}
	existing, _ := cache.FetchReference(ref)
	r.Exists = existing.Exists

	// Save the model config from the model.Model.
	config, _, err := cache.saveModelConfig(m)
	if err != nil {
		return &r, err
	}
	r.Config = config

	// Save the model content from the model.Model.
	contentLayer, _, err := cache.saveModelContentLayer(m)
	if err != nil {
		return &r, err
	}
	r.ContentLayer = contentLayer
	info, err := cache.ociStore.Info(ctx.Context(cache.out, cache.debug), contentLayer.Digest)
	if err != nil {
		return &r, err
	}
	r.Size = info.Size
	r.Digest = info.Digest
	r.CreatedAt = info.CreatedAt

	// Save the manifest for the given layers and config.
	// We do not save in index.json here, just save in blobs.
	manifest, _, err := cache.saveModelManifest(config, contentLayer)
	if err != nil {
		return &r, err
	}
	r.Manifest = manifest
	return &r, nil
}

// DeleteReference deletes a chart ref from cache
// TODO: garbage collection, only manifest removed
func (cache *Cache) DeleteReference(ref *oci.Reference) (*CacheRefSummary, error) {
	if err := cache.init(); err != nil {
		return nil, err
	}
	r, err := cache.FetchReference(ref)
	if err != nil || !r.Exists {
		return r, err
	}
	cache.ociStore.DeleteReference(r.Name)
	err = cache.ociStore.SaveIndex()
	return r, err
}

// ListReferences lists all chart refs in a cache
func (cache *Cache) ListReferences() ([]*CacheRefSummary, error) {
	if err := cache.init(); err != nil {
		return nil, err
	}
	var rr []*CacheRefSummary
	for _, desc := range cache.ociStore.ListReferences() {
		name := desc.Annotations[ocispec.AnnotationRefName]
		if name == "" {
			if cache.debug {
				fmt.Fprintf(cache.out, "warning: found manifest without name: %s", desc.Digest.Hex())
			}
			continue
		}
		ref, err := oci.ParseReference(name)
		if err != nil {
			return rr, err
		}
		r, err := cache.FetchReference(ref)
		if err != nil {
			return rr, err
		}
		rr = append(rr, r)
	}
	return rr, nil
}

// AddManifest provides a manifest to the cache index.json.
func (cache *Cache) AddManifest(ref *oci.Reference, manifest *ocispec.Descriptor) error {
	if err := cache.init(); err != nil {
		return err
	}
	cache.ociStore.AddReference(ref.FullName(), *manifest)
	err := cache.ociStore.SaveIndex()
	return err
}

// Provider provides a valid containerd Provider
func (cache *Cache) Provider() content.Provider {
	return content.Provider(cache.ociStore)
}

// Ingester provides a valid containerd Ingester
func (cache *Cache) Ingester() content.Ingester {
	return content.Ingester(cache.ociStore)
}

// ProvideIngester provides a valid oras ProvideIngester
func (cache *Cache) ProvideIngester() orascontent.ProvideIngester {
	return orascontent.ProvideIngester(cache.ociStore)
}

// init creates files needed necessary for OCI layout store
func (cache *Cache) init() error {
	if cache.ociStore == nil {
		ociStore, err := orascontent.NewOCIStore(cache.rootDir)
		if err != nil {
			return err
		}
		cache.ociStore = ociStore
		cache.memoryStore = orascontent.NewMemoryStore()
	}
	return nil
}

// fetchBlob retrieves a blob from filesystem
func (cache *Cache) fetchBlob(desc *ocispec.Descriptor) ([]byte, error) {
	reader, err := cache.ociStore.ReaderAt(ctx.Context(cache.out, cache.debug), *desc)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, desc.Size)
	_, err = reader.ReadAt(bytes, 0)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// storeBlob stores a blob on filesystem
func (cache *Cache) storeBlob(blobBytes []byte) (bool, error) {
	var exists bool
	writer, err := cache.ociStore.Store.Writer(ctx.Context(cache.out, cache.debug),
		content.WithRef(digest.FromBytes(blobBytes).Hex()))
	if err != nil {
		return exists, err
	}
	_, err = writer.Write(blobBytes)
	if err != nil {
		return exists, err
	}
	err = writer.Commit(ctx.Context(cache.out, cache.debug), 0, writer.Digest())
	if err != nil {
		if !errdefs.IsAlreadyExists(err) {
			return exists, err
		}
		exists = true
	}
	err = writer.Close()
	return exists, err
}

// saveModelConfig stores the model config as json blob and returns a descriptor
func (cache *Cache) saveModelConfig(m *model.Model) (*ocispec.Descriptor, bool, error) {
	configBytes, err := json.Marshal(m.Metadata)
	if err != nil {
		return nil, false, err
	}
	configExists, err := cache.storeBlob(configBytes)
	if err != nil {
		return nil, configExists, err
	}
	descriptor := cache.memoryStore.Add("", consts.MediaTypeModelConfig, configBytes)
	return &descriptor, configExists, nil
}

// saveModelContentLayer stores the model as tarball blob and returns a descriptor
func (cache *Cache) saveModelContentLayer(m *model.Model) (*ocispec.Descriptor, bool, error) {
	destDir := filepath.Join(cache.rootDir, ".build")
	os.MkdirAll(destDir, 0755)
	// TODO: Save models instead of letting users do it.

	contentExists, err := cache.storeBlob(m.Content)
	if err != nil {
		return nil, contentExists, err
	}
	descriptor := cache.memoryStore.Add("", consts.MediaTypeModelContentLayer, m.Content)
	return &descriptor, contentExists, nil
}

// saveModelManifest stores the image manifest as json blob and returns a descriptor
func (cache *Cache) saveModelManifest(config *ocispec.Descriptor, contentLayer *ocispec.Descriptor) (*ocispec.Descriptor, bool, error) {
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config:    *config,
		Layers:    []ocispec.Descriptor{*contentLayer},
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		return nil, false, err
	}
	manifestExists, err := cache.storeBlob(manifestBytes)
	if err != nil {
		return nil, manifestExists, err
	}
	descriptor := ocispec.Descriptor{
		MediaType: ocispec.MediaTypeImageManifest,
		Digest:    digest.FromBytes(manifestBytes),
		Size:      int64(len(manifestBytes)),
	}
	return &descriptor, manifestExists, nil
}
