package model

import "time"

type Model struct {
	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata"`
	Path     string    `json:"path"`
	// TODO(gaocegege): Tar it in the code.
	Content []byte `json:"content"`
}

type Metadata struct {
	Author          string            `json:"author,omitempty"`
	Created         time.Time         `json:"created,omitempty"`
	Description     string            `json:"description,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Format          string            `json:"format"`
	Metrics         []Metric          `json:"metrics,omitempty"`
	Hyperparameters []Hyperparameter  `json:"hyperparameters,omitempty"`
}

type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Hyperparameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
