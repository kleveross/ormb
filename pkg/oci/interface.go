package oci

import "github.com/caicloud/ormb/pkg/model"

// Interface is the interface of the client.
type Interface interface {
	Login(hostname string, username string, password string, insecure bool) error
	Logout(hostname string) error
	SaveModel(ch *model.Model, ref *Reference) error
	PushModel(ref *Reference) error
	RemoveModel(ref *Reference) error
	PullModel(ref *Reference) error
	LoadModel(ref *Reference) (*model.Model, error)
}
