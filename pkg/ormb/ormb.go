package ormb

import (
	"path/filepath"

	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
)

// ORMB is the interface to save/pull/push/export
// models in/to a remote registry.
type ORMB interface {
	Login(hostname, username, password string, insecureOpt bool) error
	Push(refStr string) error
	Pull(refStr string) error
	Export(refStr, dst string) error
	Save(src, refStr string) error
	Remove(refStr string) error
}

type ociORMB struct {
	client *oci.Client
}

// NewOCIORMB creates a OCI-based ORMB client.
func NewOCIORMB(opts ...oci.ClientOption) (ORMB, error) {
	c, err := oci.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &ociORMB{
		client: c,
	}, nil
}

func (o ociORMB) Login(hostname, username, password string, insecureOpt bool) error {
	return o.client.Login(hostname, username, password, insecureOpt)
}

func (o ociORMB) Push(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}
	return o.client.PushModel(ref)
}

func (o ociORMB) Pull(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}
	return o.client.PullModel(ref)
}

func (o ociORMB) Export(refStr, dst string) error {
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

	saver := model.NewDefaultSaver()
	if _, err := saver.Save(m, path); err != nil {
		return err
	}
	return nil
}

func (o ociORMB) Save(src, refStr string) error {
	path, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}

	l := model.NewDefaultLoader()
	m, err := l.Load(path)
	if err != nil {
		return err
	}
	return o.client.SaveModel(m, ref)
}

func (o ociORMB) Remove(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}

	return o.client.RemoveModel(ref)
}
