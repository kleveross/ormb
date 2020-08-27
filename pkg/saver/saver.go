package saver

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

	"github.com/kleveross/ormb/pkg/consts"
	"github.com/kleveross/ormb/pkg/model"
	"github.com/kleveross/ormb/pkg/parser"
)

// Saver is the implementation.
type Saver struct {
	Parser parser.Parser
}

// New creates a new Saver.
func New() Interface {
	return &Saver{
		Parser: parser.NewDefaultParser(),
	}
}

// Save saves the model from the path to the memory.
func (d Saver) Save(path string) (*model.Model, error) {
	// Save model config from <path>/ormbfile.yaml.
	dat, err := ioutil.ReadFile(filepath.Join(path, consts.ORMBfileName))
	if err != nil {
		return nil, err
	}

	metadata := &model.Metadata{}
	if metadata, err = d.Parser.Parse(dat); err != nil {
		return nil, err
	}

	format := model.Format(metadata.Format)
	if err := format.ValidateDirectory(path); err != nil {
		return nil, err
	}

	// Save the model from <path>/model.
	buf := &bytes.Buffer{}
	directoryStructure, err := TarAndGetDirectoryStructure(
		filepath.Join(path, consts.ORMBModelDirectory), buf)
	if err != nil {
		return nil, err
	}
	// Set directoryStructure for the model metadata.
	metadata.DirectoryStureture = directoryStructure

	m := &model.Model{
		Metadata: metadata,
		Path:     path,
		Config:   dat,
		Content:  buf.Bytes(),
	}
	return m, nil
}

// TarAndGetDirectoryStructure is copied from https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07.
func TarAndGetDirectoryStructure(
	src string, writers ...io.Writer) ([]string, error) {
	structure := make([]string, 0)

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return nil, fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	err := filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

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
		header.Name = strings.TrimPrefix(
			strings.Replace(file, parentDir, "", -1), string(filepath.Separator))

		// Add filename to the directory structure.
		structure = append(structure, header.Name)

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

	return structure, err
}
