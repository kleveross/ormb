package cache

import (
	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
	"github.com/containerd/containerd/content"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	orascontent "github.com/deislabs/oras/pkg/content"
)

// Interface is the interface of the cache.
type Interface interface {
	FetchReference(ref *oci.Reference) (*CacheRefSummary, error)
	StoreReference(ref *oci.Reference, m *model.Model) (*CacheRefSummary, error)
	DeleteReference(ref *oci.Reference) (*CacheRefSummary, error)
	ListReferences() ([]*CacheRefSummary, error)
	AddManifest(ref *oci.Reference, manifest *ocispec.Descriptor) error
	Provider() content.Provider
	Ingester() content.Ingester
	ProvideIngester() orascontent.ProvideIngester
}