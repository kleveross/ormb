package oras

import (
	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/oci"
)

// Interface is the interface of the client.
type Interface interface {
	Login(hostname string, username string, password string, insecure bool) error
	Logout(hostname string) error
	SaveModel(ch *model.Model, ref *oci.Reference) error
	PushModel(ref *oci.Reference) error
	RemoveModel(ref *oci.Reference) error
	PullModel(ref *oci.Reference) error
	LoadModel(ref *oci.Reference) (*model.Model, error)
	TagModel(ref *oci.Reference, target *oci.Reference) error
	Models() error
}
