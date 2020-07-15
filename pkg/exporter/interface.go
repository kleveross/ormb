package exporter

import "github.com/kleveross/ormb/pkg/model"

// Interface exports the model to the destination.
type Interface interface {
	Export(m *model.Model, dst string) (string, error)
}
