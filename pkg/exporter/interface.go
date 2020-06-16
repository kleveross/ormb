package exporter

import "github.com/caicloud/ormb/pkg/model"

// Interface exports the model to the destination.
type Interface interface {
	Export(m *model.Model, dst string) (string, error)
}
