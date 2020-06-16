package saver

import "github.com/caicloud/ormb/pkg/model"

// Interface saves the model from the path to the memory.
type Interface interface {
	Save(path string) (*model.Model, error)
} 
