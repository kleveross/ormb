package model

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/caicloud/ormb/pkg/consts"
)

// Loader loads the model from the path.
type Loader interface {
	Load(path string) (*Model, error)
}

type defaultLoader struct{}

// NewDefaultLoader creates a new defaultLoader.
func NewDefaultLoader() Loader {
	return &defaultLoader{}
}

// Load loads the model from the path.
func (d defaultLoader) Load(path string) (*Model, error) {
	// Load model config from <path>/ormbfile.yaml.
	dat, err := ioutil.ReadFile(filepath.Join(path, consts.ORMBfileName))
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{}
	if err := yaml.Unmarshal(dat, &metadata); err != nil {
		return nil, err
	}

	// Load the model from <path>/model.
	buf := &bytes.Buffer{}
	if err := Tar(filepath.Join(path, consts.ORMBModelDirectory), buf); err != nil {
		return nil, err
	}

	m := &Model{
		Metadata: metadata,
		Path:     path,
		Content:  buf.Bytes(),
	}
	return m, nil
}

// Tar is copied from https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07.
func Tar(src string, writers ...io.Writer) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files (thanks to [kumo](https://medium.com/@komuw/just-like-you-did-fbdd7df829d3) for this suggested update)
		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		parentDir := filepath.Dir(src)

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, parentDir, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}
