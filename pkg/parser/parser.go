package parser

import (
	"github.com/caicloud/ormb/pkg/model"
	"gopkg.in/yaml.v2"
)

// Parser is the type to parse the config bytes to metadata.
type Parser interface {
	Parse(configBytes []byte) (*model.Metadata, error)
}

type defaultParser struct{}

// NewDefaultParser creates a new defaultParser.
func NewDefaultParser() Parser {
	return &defaultParser{}
}

func (d defaultParser) Parse(configBytes []byte) (*model.Metadata, error) {
	metadata := &model.Metadata{}
	if err := yaml.Unmarshal(configBytes, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}
