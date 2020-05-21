package ormb

import (
	"os"
	"path/filepath"

	"github.com/caicloud/ormb/pkg/model"
	"github.com/caicloud/ormb/pkg/oci"
)

type ormb interface {
	Login(hostname, username, password string, insecureOpt bool) error
	Push(refStr string) error
	Pull(refStr string) error
	Export(refStr, dst string) error
	Save(src, refStr string) error
}

type ociormb struct {
	client *oci.Client
}

func NewDefaultOCIormb() (ormb, error) {
	c, err := oci.NewClient(oci.ClientOptWriter(os.Stdout))
	if err != nil {
		return nil, err
	}
	return &ociormb{
		client: c,
	}, nil
}

func (o ociormb) Login(hostname, username, password string, insecureOpt bool) error {
	return o.client.Login(hostname, username, password, insecureOpt)
}

func (o ociormb) Push(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		panic(err)
	}
	return o.client.PushModel(ref)
}

func (o ociormb) Pull(refStr string) error {
	ref, err := oci.ParseReference(refStr)
	if err != nil {
		return err
	}
	return o.client.PullModel(ref)
}

func (o ociormb) Export(refStr, dst string) error {
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

func (o ociormb) Save(src, refStr string) error {
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
