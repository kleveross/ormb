package exporter

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caicloud/ormb/pkg/consts"
	"github.com/caicloud/ormb/pkg/model"
	"gopkg.in/yaml.v2"
)

// Exporter exports the model to the destination.
type Exporter interface {
	Export(m *model.Model, dst string) (string, error)
}

type defaultExporter struct{}

// NewDefaultExporter creates a new defaultExporter.
func NewDefaultExporter() Exporter {
	return &defaultExporter{}
}

// Export saves the model to the destination.
func (d defaultExporter) Export(m *model.Model, dst string) (string, error) {
	if err := d.exportMetadata(m, dst); err != nil {
		return "", err
	}

	// Export model.
	gzr, err := gzip.NewReader(bytes.NewBuffer(m.Content))
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return "", nil

		// return any other error
		case err != nil:
			return "", err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return "", err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			// Create the parent directory first.
			parentDir := filepath.Dir(target)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return "", err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return "", err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func (d defaultExporter) exportMetadata(m *model.Model, dst string) error {
	// Export ormbfile.yaml.
	yamlBytes, err := yaml.Marshal(m.Metadata)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dst, consts.ORMBfileName), yamlBytes, 0644)
}
