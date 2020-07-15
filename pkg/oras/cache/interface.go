package cache

import (
	"github.com/containerd/containerd/content"
	orascontent "github.com/deislabs/oras/pkg/content"
	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/oci"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
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
