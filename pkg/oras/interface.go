package oras

import (
	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
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
}
