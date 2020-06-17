package ormb

import (
	"path/filepath"

	"github.com/caicloud/ormb/pkg/exporter"
	"github.com/caicloud/ormb/pkg/oci"
	"github.com/caicloud/ormb/pkg/oras"
	"github.com/caicloud/ormb/pkg/saver"
)

// Interface is the interface to manage
// models with a remote registry.
type Interface interface {
	Login(hostname, username, password string, insecureOpt bool) error
	Push(refStr string) error
	Pull(refStr string) error
	Export(refStr, dst string) error
	Save(src, refStr string) error
	Remove(refStr string) error
}

type ORMB struct {
	client   oras.Interface
	saver    saver.Interface
	exporter exporter.Interface
}

// New creates a OCI-based ORMB client.
func New(opts ...oras.ClientOption) (Interface, error) {
	c, err := oras.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &ORMB{
		client:   c,
		saver:    saver.New(),
		exporter: exporter.New(),
	}, nil
}

func (o ORMB) Login(hostname, username, password string, insecureOpt bool) error {
	return o.client.Login(hostname, username, password, insecureOpt)
}

func (o ORMB) Push(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}
	return o.client.PushModel(ref)
}

func (o ORMB) Pull(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}
	return o.client.PullModel(ref)
}

func (o ORMB) Export(refStr, dst string) error {
	path, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}

	m, err := o.client.LoadModel(ref)
	if err != nil {
		return err
	}

	if _, err := o.exporter.Export(m, path); err != nil {
		return err
	}
	return nil
}

func (o ORMB) Save(src, refStr string) error {
	path, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}

	m, err := o.saver.Save(path)
	if err != nil {
		return err
	}
	return o.client.SaveModel(m, ref)
}

func (o ORMB) Remove(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}

	return o.client.RemoveModel(ref)
}
